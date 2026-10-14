package queue

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"vio/models"

	"github.com/resend/resend-go/v4"
)

var taskChannel chan models.Task = make(chan models.Task, 10)
var logChannel chan models.Log = make(chan models.Log, 10)

func sendEmail(text string) error {
	client := resend.NewClient(os.Getenv("RESEND_API_KEY"))
	params := &resend.SendEmailRequest{
		To:   []string{os.Getenv("RECIPIENT_EMAIL")},
		From: os.Getenv("MAIL_FROM"),
		Html: text,
		Text: text,
	}

	_, err := client.Emails.Send(params)

	return err
}

func processTasks(id int) {
	for task := range taskChannel {
		start := time.Now()
		var end time.Time
		var err error
		var status int

		if task.Type == 1 {
			err = sendEmail(task.Content)
			end = time.Now()
		} else {
			fmt.Println("This is a worker for image generation")
		}

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

		logChannel <- models.Log{Type: task.Type, Content: task.Content, Status: status, Duration: float64(duration.Milliseconds()), Timestamp: time.Now(), Error: error}
	}
}

func processLogs() {
	for log := range logChannel {
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

func Queue(task models.Task) {
	taskChannel <- task

	initProcessing := func() {
		for i := 1; i <= 3; i++ {
			go processTasks(i)
		}
		go processLogs()
		fmt.Println("Inside the init processing method")
	}

	if (len(taskChannel) - 1) == 0 {
		initProcessing()
	}
}
