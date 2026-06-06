import (
	"fmt"
	"sync"
)

func main() {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

	var wg sync.WaitGroup
	evench := make(chan int, len(numbers))
	oddch := make(chan int, len(numbers))
	mulfivech := make(chan int, len(numbers))

	wg.Add(3)
	go func() {
		defer wg.Done()
		for number := range evench {
			if number%2 == 0 {
				fmt.Printf("Even: %d\n", number)
			}
		}
	}()

	go func() {
		defer wg.Done()
		for number := range oddch {
			if number%2 != 0 {
				fmt.Printf("Odd: %d\n", number)
			}
		}
	}()

	go func() {
		defer wg.Done()
		for number := range mulfivech {
			if number%5 == 0 {
				fmt.Printf("Five Multiple: %d\n", number)
			}
		}
	}()

	for _, number := range numbers {
		evench <- number
		oddch <- number
		mulfivech <- number
	}

	close(evench)
	close(oddch)
	close(mulfivech)
	wg.Wait()
}
