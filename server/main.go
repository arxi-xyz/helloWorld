package main

import (
	"fmt"
	"net/http"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 10000; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			resp, err := http.Get("http://localhost:8111/add")
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			defer resp.Body.Close()

			fmt.Println(resp.StatusCode)
		}()
	}

	wg.Wait()
}
