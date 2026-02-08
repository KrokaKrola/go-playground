package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	chunks := []int{1, 2, 3}

	fmt.Println("Manager: Starting batch job...")

	for _, id := range chunks {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			processChunk(ctx, workerID, cancel)
		}(id)
	}

	// waiting for everyone to finish (either by success or cancelation)
	wg.Wait()
	fmt.Println("Manager: All workers have shut down.")
}

func processChunk(ctx context.Context, workerID int, cancel context.CancelFunc) {
	fmt.Printf("Worker %d: Starting...\n", workerID)

	if workerID == 2 {
		time.Sleep(time.Duration(workerID) * time.Second)
		fmt.Printf("Worker %d: CRITICAL ERROR. Data corrupt. Stopping whole pipeline\n", workerID)
		// trigger cancelation for all related goroutines
		cancel()
		return
	}

	select {
	case <-time.After(time.Duration(workerID) * time.Second):
		fmt.Printf("Worker %d: Finished successfully.\n", workerID)
	case <-ctx.Done():
		fmt.Printf("Worker %d: Cancel was called from one of the workers. Qutting early\n", workerID)
		return
	}
}
