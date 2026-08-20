package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Url struct {
	i   int
	url string
}

type downLoadedData struct {
	i       int
	content string
}

func downloadOne(ctx context.Context, url Url) (downLoadedData, error) {
	select {
	case <-time.After(2 * time.Second):
	case <-ctx.Done():
		fmt.Printf("\nDownload aborted %d", url.i)
		return downLoadedData{}, ctx.Err()
	}

	if url.i == 3 {
		return downLoadedData{}, fmt.Errorf("File not found %d", 3)
	}

	fmt.Println("Download done: " + url.url)
	return downLoadedData{i: url.i, content: "Downloaded" + url.url}, nil
}

func main() {
	urls := []string{"abc", "pqr", "aa", "sdsd", "def"}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	resultChan := make(chan downLoadedData, len(urls))
	errChan := make(chan error, len(urls))

	// 1. Launch downloads concurrently
	for i, url := range urls {
		wg.Add(1)

		go func(i int, url string) {
			defer wg.Done()

			result, err := downloadOne(ctx, Url{i: i, url: url})
			if err != nil {
				errChan <- err
				cancel() // Trigger cancellation for peer downloads
				return
			}

			// Unconditionally send result to avoid losing completed work
			resultChan <- result
		}(i, url)
	}

	// 2. Offload wg.Wait() to a background goroutine so main doesn't block.
	// This ensures main can stream results immediately as they arrive.
	go func() {
		wg.Wait()
		close(resultChan)
		close(errChan)
	}()

	downLoadedDataSlice := make([]string, len(urls))

	// 3. Main streams incoming results in real-time concurrently
	for result := range resultChan {
		fmt.Println("Received result:", result)
		downLoadedDataSlice[result.i] = result.content
	}

	// 4. Safely drain all collected errors once resultChan closes
	for err := range errChan {
		fmt.Println("Error:", err)
	}

	fmt.Println("Final Slice Result:", downLoadedDataSlice)
}
