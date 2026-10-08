package queue

import (
	"sync"
	"vio/models"
)

var taskChannel chan models.Task = make(chan models.Task, 10)
var resultChannel chan models.Result = make(chan models.Result, 10)

var taskWg sync.WaitGroup
var resultWg sync.WaitGroup

func processTask(tasks <-chan models.Task, results chan<- models.Result) {

}

func queue(task models.Task) {
	_, isTaskChannelClosed := <-taskChannel

	if isTaskChannelClosed {
		taskChannel = make(chan models.Task, 10)
	}

	if _, isResultChannelClosed := <-resultChannel; isResultChannelClosed {
		resultChannel = make(chan models.Result, 10)
	}

	taskWg.Add(1)
	resultWg.Add(1)

	initTaskProcessing := func() {
		for i := 1; i <= 3; i++ {
			go processTask(taskChannel, resultChannel)
		}
	}

	// if isTaskChannelClosed ||
}
