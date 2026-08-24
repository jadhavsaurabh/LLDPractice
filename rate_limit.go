package main

import (
	"fmt"
	"time"
)

func main() {
	reqCh := make(chan int, 20)

  // Producer routine
	go func() {
		defer close(reqCh)
		for num := range 20 {
			reqCh <- num
		}
	}()

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	// Direct channel synchronization
  // chanel wait karvayega
	for req := range reqCh {
		<-ticker.C
		fmt.Printf("[%s] Processing req %d\n", time.Now().Format("15:04:05.000"), req)
	}
}
