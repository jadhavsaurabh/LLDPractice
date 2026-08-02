package main

import (
	"fmt"
	"sync"
	"time"
)

type UrlCounter struct {
	mp map[string]int
	mu sync.RWMutex
}

func (urlC *UrlCounter) hit(url string) {
	urlC.mu.Lock()
	defer urlC.mu.Unlock()

	fmt.Printf("\n hitting %s url", url)
	urlC.mp[url]++
}

func (urlC *UrlCounter) read() {
	urlC.mu.RLock()
	defer urlC.mu.RUnlock()

	fmt.Printf("\n Map: %+v", urlC.mp)
}

func main() {
	urlC := &UrlCounter{
		mp: make(map[string]int),
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		urls := []string{"abc", "pqr", "abc", "abc", "pqr", "xyz", "xyz", "abc", "pqr", "abc", "pqr", "abc", "abc", "pqr", "xyz", "xyz", "abc", "pqr", "abc", "pqr", "abc", "abc", "pqr", "xyz", "xyz", "abc", "pqr", "abc", "pqr", "abc", "abc", "pqr", "xyz", "xyz", "abc", "pqr", "abc", "pqr", "abc", "abc", "pqr", "xyz", "xyz", "abc", "pqr", "abc", "pqr", "abc", "abc", "pqr", "xyz", "xyz", "abc", "pqr", "abc", "pqr", "abc", "abc", "pqr", "xyz", "xyz", "abc", "pqr", "abc", "pqr", "abc", "abc", "pqr", "xyz", "xyz", "abc", "pqr", "abc", "pqr", "abc", "abc", "pqr", "xyz", "xyz", "abc", "pqr", "abc", "pqr", "abc", "abc", "pqr", "xyz", "xyz", "abc", "pqr"}

		for _, url := range urls {
			urlC.hit(url)
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		defer wg.Done()
		for range 5 {
			time.Sleep(8 * time.Second)
			urlC.read()
		}
	}()

	wg.Wait()
}
