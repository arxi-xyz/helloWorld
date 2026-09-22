package main

import (
	"fmt"
	"sync"
	"testing"
)

type lockerCache interface {
	Get(key string) (string, bool)
	Set(key, value string)
}

type mutexCache struct {
	mu   sync.Mutex
	data map[string]string
}

func newMutexCache() *mutexCache {
	return &mutexCache{data: make(map[string]string)}
}

func (c *mutexCache) Get(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, ok := c.data[key]
	return value, ok
}

func (c *mutexCache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

type rwMutexCache struct {
	mu   sync.RWMutex
	data map[string]string
}

func newRWMutexCache() *rwMutexCache {
	return &rwMutexCache{data: make(map[string]string)}
}

func (c *rwMutexCache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.data[key]
	return value, ok
}

func (c *rwMutexCache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

func benchmarkCache(b *testing.B, cache lockerCache, readPercent int) {
	const nkeys = 256
	keys := make([]string, nkeys)
	for i := range keys {
		keys[i] = fmt.Sprintf("k%d", i)
		cache.Set(keys[i], "v")
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := keys[i%nkeys]
			if i%100 < readPercent {
				cache.Get(key)
			} else {
				cache.Set(key, "v")
			}
			i++
		}
	})
}

func BenchmarkCacheLock(b *testing.B) {
	workloads := []struct {
		name        string
		readPercent int
	}{
		{name: "Read100", readPercent: 100},
		{name: "Read90", readPercent: 90},
		{name: "Read50", readPercent: 50},
		{name: "Write100", readPercent: 0},
	}

	for _, wl := range workloads {
		b.Run("Mutex/"+wl.name, func(b *testing.B) {
			benchmarkCache(b, newMutexCache(), wl.readPercent)
		})
		b.Run("RWMutex/"+wl.name, func(b *testing.B) {
			benchmarkCache(b, newRWMutexCache(), wl.readPercent)
		})
	}
}
