package main

import (
	"fmt"
	"sync"
	"time"
)

type Counter struct {
	mu  sync.RWMutex
	val int
}

// Increment safely increments the counter.
func (c *Counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()

	fmt.Printf("\n[Writer] Current Value: %d", c.val)
	c.val++
	fmt.Printf("\n[Writer] New Value    : %d", c.val)
}

// Value safely returns the current counter value.
func (c *Counter) Value() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.val
}

func main() {
	var wg sync.WaitGroup

	counter := &Counter{}

	wg.Add(4)

	// Writer
	go func() {
		defer wg.Done()

		for i := 0; i < 5; i++ {
			counter.Increment()
			time.Sleep(300 * time.Millisecond)
		}
	}()

	// Reader 1
	go func() {
		defer wg.Done()

		for i := 0; i < 5; i++ {
			fmt.Printf("\n[Reader 1] Value: %d", counter.Value())
			time.Sleep(150 * time.Millisecond)
		}
	}()

	// Reader 2
	go func() {
		defer wg.Done()

		for i := 0; i < 5; i++ {
			fmt.Printf("\n[Reader 2] Value: %d", counter.Value())
			time.Sleep(200 * time.Millisecond)
		}
	}()

	// Reader 3
	go func() {
		defer wg.Done()

		for i := 0; i < 5; i++ {
			fmt.Printf("\n[Reader 3] Value: %d", counter.Value())
			time.Sleep(250 * time.Millisecond)
		}
	}()

	wg.Wait()

	fmt.Printf("\n\nFinal Counter Value: %d\n", counter.Value())
}
