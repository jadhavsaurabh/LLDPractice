/*

Requirements:

Download 100 files concurrently
Limit max concurrency to 5
Retry failed downloads
Timeout support using context.Context

*/

import (
	"context"
	"sync"
	"time"
	"fmt"
)

func Download(ctx context.Context, url string) ([]byte, error) {
    fmt.Println("Downloading", url)

    select {
    case <-time.After(3 * time.Second):
        fmt.Println("Downloaded", url)
        return []byte{}, nil

    case <-ctx.Done():
        return nil, ctx.Err()
    }
}

func worker(
    ctx context.Context,
    jobs <- chan string,
    wg *sync.WaitGroup,
) {
    defer wg.Done()

    for url := range jobs {
        ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)

        _, err := Download(ctxTimeout, url)
        if err != nil {
            fmt.Println("failed:", url, err)
        }

        cancel()
    }
}

func main() {
    urls := []string{
        "1", "2", "3", "4", "5",
        "6", "7", "8", "9", "10",
    }

    jobs := make(chan string)

    var wg sync.WaitGroup

    workerCount := 5

    ctx := context.Background()

    for i := 0; i < workerCount; i++ {
        wg.Add(1)
        go worker(ctx, jobs, &wg)
    }

    for _, url := range urls {
        jobs <- url
    }

    close(jobs)

    wg.Wait()
}
