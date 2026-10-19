package queue

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"time"
	"vio/models"

	"github.com/resend/resend-go/v4"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
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

func generateImage(text string) error {
	img := image.NewRGBA(image.Rect(0, 0, 400, 150))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{30, 40, 50, 255}}, image.Point{}, draw.Src)

	parsedFont, err := opentype.Parse(goregular.TTF)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to parse font: %v", err))
	}

	face, err := opentype.NewFace(parsedFont, &opentype.FaceOptions{Size: 28, DPI: 72, Hinting: font.HintingNone})
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to create font face: %v", err))
	}

	drawer := &font.Drawer{Dst: img, Src: image.NewUniform(color.RGBA{255, 255, 255, 255}), Face: face, Dot: fixed.Point26_6{X: fixed.I(50), Y: fixed.I(85)}}
	drawer.DrawString(text)

	directory, _ := os.Getwd()
	imagePath := filepath.Join(directory, "images", fmt.Sprintf("%s.png", text))
	file, err := os.Create(imagePath)
	if err != nil {
		return errors.New(fmt.Sprintf("There was a problem creating the image: %v", err))
	}

	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		return errors.New(fmt.Sprintf("Failed to encode PNG: %v", err))
	}

	return nil
}

func processTasks(id int) {
	for task := range taskChannel {
		start := time.Now()
		var end time.Time
		var err error
		var status int

		if task.Type == 1 {
			err = sendEmail(task.Content)
		} else {
			err = generateImage(task.Content)
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
