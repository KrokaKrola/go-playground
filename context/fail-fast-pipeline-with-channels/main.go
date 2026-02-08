package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resultsCh := make(chan int)
	var wg sync.WaitGroup

	dataChunks := []int{10, 20, 30, 40, 50}

	fmt.Println("Manager: Starting workers...")

	for _, chunk := range dataChunks {
		wg.Add(1)

		go func(val int) {
			defer wg.Done()
			processChunk(ctx, val, resultsCh, cancel)
		}(chunk)
	}

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	finalResults := []int{}

	for res := range resultsCh {
		fmt.Printf("Got data from results channel %d\n", res)
		finalResults = append(finalResults, res)
	}

	if ctx.Err() != nil {
		fmt.Println("Manager: Job failed/canceld. Partial results: ", finalResults)
	} else {
		fmt.Printf("Manager: Success. All results collected: %v\n", finalResults)
	}
}

func processChunk(ctx context.Context, val int, out chan<- int, cancel context.CancelFunc) {
	select {
	case <-time.After(time.Duration(val*100) * time.Millisecond):
		// simmulate failure on chunk 30
		if val == 30 {
			fmt.Printf("Worker %d: Error found. Canceling whole pipeline\n", val)
			cancel()
			return
		}

		select {
		case out <- val * 2:
			// sent data to parent
		case <-ctx.Done():
			// context was canceled while trying to send data
			return
		}
	case <-ctx.Done():
		return
	}
}
