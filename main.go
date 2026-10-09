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
	queueMux := http.NewServeMux()
	requestMux := http.NewServeMux()

	controllers.RegisterQueueRoutes(queueMux)
	controllers.RegisterRequestRoutes(requestMux)

	ctx, cancelCtx := context.WithCancel(context.Background())
	queueServer := &http.Server{
		Addr:    ":8000",
		Handler: queueMux,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}

	requestServer := &http.Server{
		Addr:    ":5000",
		Handler: requestMux,
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
		}

		cancelCtx()
	}()

	go func() {
		err := requestServer.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Queue server port closed")
		} else if err != nil {
			fmt.Printf("There was a problem running the request server: %s", err)
		}

		cancelCtx()
	}()

	<-ctx.Done()
}
