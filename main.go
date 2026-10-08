package main

import (
	"errors"
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	queueServer := &http.Server{
		Addr:    ":8000",
		Handler: mux,
	}

	go func() {
		err := queueServer.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Queue server port closed")
		} else if err != nil {
			fmt.Printf("There was a problem running the queue server: %s", err)
		}
	}()
}
