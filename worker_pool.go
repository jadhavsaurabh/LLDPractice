package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type Job struct {
	ID int
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobCh := make(chan Job, 5)
	resultCh := make(chan int, 5)
	errCh := make(chan error, 1)

	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go worker(ctx, cancel, &wg, jobCh, resultCh, errCh)
	}

	// Producer
	go produce(ctx, jobCh)

	// Close result channel after all workers exit
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Stream results as soon as workers produce them
	for result := range resultCh {
		fmt.Println("result:", result)
	}

	select {
	case err := <-errCh:
		fmt.Println("error:", err)
	default:
	}
}

func worker(
	ctx context.Context,
	cancel context.CancelFunc,
	wg *sync.WaitGroup,
	jobCh <-chan Job,
	resultCh chan<- int,
	errCh chan<- error,
) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return

		case job, ok := <-jobCh:
			if !ok {
				return
			}

			// Simulate an error
			if job.ID == 7 {
				select {
				case errCh <- errors.New("failed to process job"):
				default:
				}

				cancel()
				return
			}

			result := job.ID * job.ID

			select {
			case resultCh <- result:
			case <-ctx.Done():
				return
			}
		}
	}
}

func produce(ctx context.Context, jobCh chan<- Job) {
	for i := 0; i < 100; i++ {
		select {
		case jobCh <- Job{ID: i}:
		case <-ctx.Done():
			close(jobCh)
			return
		}
	}

	close(jobCh)
}
