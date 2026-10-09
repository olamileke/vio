package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"vio/controllers"
)

func main() {
	mux := http.NewServeMux()

	controllers.RegisterQueueRoutes(mux)

	ctx, cancelCtx := context.WithCancel(context.Background())
	queueServer := &http.Server{
		Addr:    ":8000",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		err := queueServer.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Queue server port closed")
		} else if err != nil {
			fmt.Printf("There was a problem running the queue server: %s", err)
		} else {
			fmt.Printf("Queue server running on port :8000")
		}

		cancelCtx()
	}()

	<-ctx.Done()
}
