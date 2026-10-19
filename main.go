package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"vio/controllers"
	"vio/models"
	queuepkg "vio/queue"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	queueMux := http.NewServeMux()
	requestMux := http.NewServeMux()

	var wg sync.WaitGroup
	taskChannel := make(chan models.Task, 10)
	logChannel := make(chan models.Log, 10)

	queue := queuepkg.Queue{Wg: wg, TaskChannel: taskChannel, LogChannel: logChannel}

	controllers.RegisterQueueRoutes(queueMux, &queue)
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
			fmt.Printf("Queue server port closed\n")
		} else if err != nil {
			fmt.Printf("There was a problem running the request server: %s\n", err)
		}

		cancelCtx()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	fmt.Println("Server is shutting down...\n")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if err := queueServer.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("Queue server forced to shutdown: %s\n", err)
	}

	fmt.Println("Server stopped\n")
}
