package main

import (
	"fmt"
	"sync"
)

func main() {
	jobCh := make(chan int)

	var wg sync.WaitGroup
	wg.Add(2)
	go generateJobs(jobCh, &wg)
	go processJobs(jobCh, &wg)

	wg.Wait()
}

func generateJobs(ch chan int, group *sync.WaitGroup) {
	defer group.Done()
	defer close(ch)

	jobs := [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}

	for _, job := range jobs {
		fmt.Printf("\nPushing job: %d", job)
		ch <- job
	}
}

func processJobs(ch chan int, group *sync.WaitGroup) {
	defer group.Done()

	for job := range ch {
		fmt.Printf("\nProcessing job: %d", job)
	}
}
