package main

import (
	"fmt"
	"sync"
)

func main() {
	commCh := make(chan string)
	var wg sync.WaitGroup

	wg.Add(2)
	go ping(commCh, &wg)
	go pong(commCh, &wg)

	wg.Wait()
}

func ping(comCh chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(comCh)

	for range 5 {
		comCh <- "ping"
		
		msg := <-comCh
		fmt.Printf("Received %s, sending next ping\n", msg)
	}
}

func pong(comCh chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for msg := range comCh {
		fmt.Printf("Received %s, responding with pong\n", msg)
		comCh <- "pong"
	}
}
