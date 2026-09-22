package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"testing"
)

func TestGetMissingKey(t *testing.T) {
	c := NewCache()

	value, ok := c.Get("missing")
	if ok {
		t.Fatalf("Get returned ok=true for missing key, value=%q", value)
	}
	if value != "" {
		t.Fatalf("Get returned value %q for missing key", value)
	}
}

func TestSetAndGet(t *testing.T) {
	c := NewCache()

	c.Set("name", "world")

	value, ok := c.Get("name")
	if !ok {
		t.Fatal("Get returned ok=false for an existing key")
	}
	if value != "world" {
		t.Fatalf("Get returned %q, want %q", value, "world")
	}
}

func TestSetOverwritesExistingValue(t *testing.T) {
	c := NewCache()

	c.Set("name", "first")
	c.Set("name", "second")

	value, ok := c.Get("name")
	if !ok {
		t.Fatal("Get returned ok=false after overwrite")
	}
	if value != "second" {
		t.Fatalf("Get returned %q, want %q", value, "second")
	}
}

func TestDeleteRemovesKey(t *testing.T) {
	c := NewCache()

	c.Set("name", "world")
	c.Delete("name")

	value, ok := c.Get("name")
	if ok {
		t.Fatalf("Get returned ok=true after Delete, value=%q", value)
	}
}

func TestDeleteMissingKey(t *testing.T) {
	c := NewCache()

	c.Set("name", "world")
	c.Delete("missing")

	value, ok := c.Get("name")
	if !ok || value != "world" {
		t.Fatalf("Delete of a missing key changed existing data: got %q, ok=%v", value, ok)
	}
}

func TestSnapshotSortedPrintsKeysInOrder(t *testing.T) {
	c := NewCache()
	c.Set("b", "two")
	c.Set("a", "one")
	c.Set("c", "three")

	output := captureStdout(t, func() {
		c.SnapshotSorted()
	})

	want := "a: one\nb: two\nc: three\n"
	if output != want {
		t.Fatalf("SnapshotSorted output = %q, want %q", output, want)
	}
}

func TestSnapshotSortedEmpty(t *testing.T) {
	c := NewCache()

	output := captureStdout(t, func() {
		c.SnapshotSorted()
	})

	if output != "" {
		t.Fatalf("SnapshotSorted on empty cache printed %q", output)
	}
}

func TestConcurrentReadersAndWriters(t *testing.T) {
	c := NewCache()

	const goroutines = 8
	const keysPerWriter = 50

	var wg sync.WaitGroup
	wg.Add(goroutines * 2)
	start := make(chan struct{})

	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
	
			<-start
	
			for n := 0; n < keysPerWriter; n++ {
				key := fmt.Sprintf("w%d-k%d", id, n)
				c.Set(key, fmt.Sprintf("v-%d-%d", id, n))
			}
		}(i)
	}
	
	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
	
			<-start
	
			for n := 0; n < keysPerWriter; n++ {
				key := fmt.Sprintf("w%d-k%d", id, n)
				c.Get(key)
			}
		}(i)
	}
	
	close(start)
	wg.Wait()

	for i := 0; i < goroutines; i++ {
		for n := 0; n < keysPerWriter; n++ {
			key := fmt.Sprintf("w%d-k%d", i, n)
			want := fmt.Sprintf("v-%d-%d", i, n)
			got, ok := c.Get(key)
			if !ok || got != want {
				t.Fatalf("Get(%q) = %q, %v; want %q, true", key, got, ok, want)
			}
		}
	}
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	original := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}

	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = original

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("reading stdout: %v", err)
	}
	return buf.String()
}
