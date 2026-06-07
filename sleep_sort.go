package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	input := [6]int{50, 40, 10, 20, 30, 11}

	var wg sync.WaitGroup
	for _, number := range input {
		wg.Add(1)
		go sleepPrint(number, &wg)
	}

	wg.Wait()
}

func sleepPrint(number int, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Second * time.Duration(number))

	fmt.Println(number)
}
