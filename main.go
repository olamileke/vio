package main

import (
	"context"
	"errors"
	"fmt"
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

	var workerWg, logWg sync.WaitGroup
	workerWg.Add(3)
	logWg.Add(1)
	taskChannel := make(chan models.Task, 10)
	logChannel := make(chan models.Log, 10)

	queue := queuepkg.Queue{WorkerWg: &workerWg, LogWg: &logWg, TaskChannel: taskChannel, LogChannel: logChannel}

	controllers.RegisterQueueRoutes(queueMux, &queue)

	server := &http.Server{
		Addr:    ":8000",
		Handler: queueMux,
	}

	go func() {
		err := server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Queue server port closed")
		} else if err != nil {
			fmt.Printf("There was a problem running the queue server: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	fmt.Println("Server is shutting down...\n")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("Queue server forced to shutdown: %s\n", err)
	}

	fmt.Println("Server stopped\n")

	queueChan := make(chan string, 1)

	go func() {
		close(queue.TaskChannel)
		close(queue.LogChannel)
		queue.WorkerWg.Wait()
		queue.LogWg.Wait()
		close(queueChan)
	}()

	select {
	case <-queueChan:
		fmt.Println("All active workers finished cleanly\n")
	case <-shutdownCtx.Done():
		fmt.Println("Shutdown timed out before all workers completed\n")
	}

	fmt.Println("Server exiting\n")
}
