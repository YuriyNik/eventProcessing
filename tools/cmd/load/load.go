package main

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
)

func main() {
	const from = 125000
	const total = 150000
	const workers = 150

	sem := make(chan struct{}, workers)
	var wg sync.WaitGroup
	fmt.Println("Starting load generator for ", total-from, "requests with", workers, "concurrent workers")
	for i := from; i < total; i++ {
		wg.Add(1)
		sem <- struct{}{}

		go func(n int) {
			defer wg.Done()
			body := []byte(fmt.Sprintf(`{"id": "%d", "payload": "hello-%d"}`, n, n))
			http.Post("http://localhost:8082/send", "application/json", bytes.NewBuffer(body))
			<-sem
		}(i)
	}
	fmt.Println("Waiting for goroutines to finish")
	wg.Wait()
	fmt.Println("All requests sent")
}
