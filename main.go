package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Counter struct {
	mu      sync.Mutex
	counter int
}

func (c *Counter) handler(w http.ResponseWriter, r *http.Request) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.counter++
}

func (c *Counter) printHandler(w http.ResponseWriter, r *http.Request) {
	c.mu.Lock()
	defer c.mu.Unlock()

	fmt.Fprintf(w, "Counter: %d", c.counter)
}

func main() {
	counter := &Counter{}

	http.HandleFunc("/add", counter.handler)
	http.HandleFunc("/print", counter.printHandler)
	log.Fatal(http.ListenAndServe(":8111", nil))
}
