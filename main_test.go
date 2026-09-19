package main

import (
	"sync"
	"testing"
)

func TestCounterConcurrent(t *testing.T) {
	counter := &Counter{}

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			counter.handler(nil, nil)
		}()
	}

	wg.Wait()

	counter.mu.Lock()
	defer counter.mu.Unlock()

	if counter.counter != 1000 {
		t.Fatalf("expected 1000, got %d", counter.counter)
	}
}
