package queue

import (
	"fmt"
	"vio/models"
)

var taskChannel chan models.Task = make(chan models.Task, 10)
var resultChannel chan models.Result = make(chan models.Result, 10)

func processTask(id int, tasks <-chan models.Task, results chan<- models.Result) {
	for task := range tasks {
		task.Content = "abcd"
		fmt.Printf("goroutine with id %d running\n", id)
	}
}

func Queue(task models.Task) {
	taskChannel <- task

	initTaskProcessing := func() {
		for i := 1; i <= 3; i++ {
			go processTask(i, taskChannel, resultChannel)
		}
		fmt.Println("Inside the init processing method")
	}

	if (len(taskChannel) - 1) == 0 {
		initTaskProcessing()
	}
}
