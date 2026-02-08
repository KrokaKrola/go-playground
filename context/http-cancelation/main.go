package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/report", reportHandler)
	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}

func reportHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handler: Request started")

	// context is automatically created by the net/http package
	// it's calceled if the user disconnects or the request times out
	ctx := r.Context()

	err := generatePDF(ctx)
	if err != nil {
		fmt.Println("Handler: " + err.Error())
		return
	}

	fmt.Fprintln(w, "PDF Generated successfully")
	fmt.Println("Handler: Request finished successfully")
}

func generatePDF(ctx context.Context) error {
	count := 5

	for i := range count {
		select {
		case <-ctx.Done():
			return fmt.Errorf("operation canceled by client")
		default:
			time.Sleep(1 * time.Second)
			fmt.Printf("Worker: Generating page %d...\n", i)
		}
	}

	return nil
}
