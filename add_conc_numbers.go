package main

import (
	"fmt"
)

func main() {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

	numParts := 4
	chunkSize := len(numbers) / numParts
	sumChan := make(chan int, numParts)

	for part := range numParts {
		go func(part int) {
			fmt.Println(part)
			start := part * chunkSize
			end := start + chunkSize
			if part == numParts-1 {
				end = len(numbers)
			}
			sum := 0
			for i := start; i < end; i++ {
				sum += numbers[i]
			}
			fmt.Println(sum)
			sumChan <- sum
		}(part)
	}

	total := 0
	for range numParts {
		total += <-sumChan
	}

	fmt.Print(total)
}
