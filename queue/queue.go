package queue

import (
	"fmt"
	"time"
	"vio/models"
)

var taskChannel chan models.Task = make(chan models.Task, 10)
var resultChannel chan models.Result = make(chan models.Result, 10)

// var taskWg sync.WaitGroup
// var resultWg sync.WaitGroup

func processTask(tasks <-chan models.Task, results chan<- models.Result) {
	fmt.Println("Inisde the goroutine")
	for task := range tasks {
		fmt.Println(task.Content)
	}
}

func Queue(task models.Task) {
	taskChannel <- task

	fmt.Printf("Here is the length of the task channel: %d\n", len(taskChannel))

	initTaskProcessing := func() {
		for i := 1; i <= 3; i++ {
			go processTask(taskChannel, resultChannel)
			time.Sleep(3 * time.Second)
		}
		fmt.Println("Inside init processing")
	}

	if (len(taskChannel) - 1) == 0 {
		initTaskProcessing()
	}
}
