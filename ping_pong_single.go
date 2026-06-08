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

	comCh <- "ping"

	msg := <-comCh
	if msg == "pong" {
		fmt.Println("Received pong closing the channel")
		close(comCh)
	}
}

func pong(comCh chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	msg := <-comCh
	if msg == "ping" {
		fmt.Println("Received ping")
		fmt.Println("Responding with pong")
		comCh <- "pong"
	}
}
