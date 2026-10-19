package queue

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
	"vio/models"
	"vio/utilities"
)

type Queue struct {
	Wg          sync.WaitGroup
	TaskChannel chan models.Task
	LogChannel  chan models.Log
}

func (queue *Queue) Process(task models.Task) {
	queue.TaskChannel <- task
	queue.Wg.Add(1)

	initProcessing := func() {
		for i := 1; i <= 3; i++ {
			go queue.processTasks(i)
		}
		go queue.processLogs()
		fmt.Println("Inside the init processing method")
	}

	if (len(queue.TaskChannel) - 1) == 0 {
		initProcessing()
	}
}

func (queue *Queue) processTasks(id int) {
	for task := range queue.TaskChannel {
		start := time.Now()
		var end time.Time
		var err error
		var status int

		if task.Type == 1 {
			err = utilities.SendEmail(task.Content)
		} else {
			err = utilities.GenerateImage(task.Content)
		}

		end = time.Now()
		duration := end.Sub(start)

		if err != nil {
			status = 2
		} else {
			status = 1
		}

		var error *string

		if err != nil {
			errorVariable := err.Error()
			error = &errorVariable
		} else {
			error = nil
		}

		queue.LogChannel <- models.Log{Type: task.Type, Content: task.Content, Status: status, Duration: float64(duration.Milliseconds()), Timestamp: time.Now(), Error: error}
	}
}

func (queue *Queue) processLogs() {
	for log := range queue.LogChannel {
		directory, _ := os.Getwd()
		logsPath := filepath.Join(directory, "logs.json")
		logsContent, _ := os.ReadFile(logsPath)

		var logsData models.LogData
		json.Unmarshal(logsContent, &logsData)

		logsData.Logs = append(logsData.Logs, log)

		updatedLogsContent, _ := json.Marshal(logsData)
		os.WriteFile(logsPath, updatedLogsContent, 0644)
	}
}
