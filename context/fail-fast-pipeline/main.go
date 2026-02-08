package main

import (
	"context"
	"fmt"
	"sync"
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

	wg.Wait()
	fmt.Println("Manager: All workers have shut down.")

}

func processChunk(ctx context.Context, workerID int, cancel context.CancelFunc) {
	panic("unimplemented")
}
