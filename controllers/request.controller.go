package controllers

import (
	"bytes"
	cryptorand "crypto/rand"
	"encoding/json"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"sync"
	"time"
	"vio/models"
)

type requestController struct{}

func RegisterRequestRoutes(mux *http.ServeMux) {
	rc := requestController{}
	mux.HandleFunc("/make-requests", rc.makeRequests)
}

func (rc *requestController) makeRequests(w http.ResponseWriter, r *http.Request) {
	taskDtos := make([]models.CreateTaskDto, 0, 15)

	for i := 1; i <= 15; i++ {
		var taskType string

		if i%3 == 0 || i%5 == 0 {
			taskType = "email"
		} else {
			taskType = "pdf"
		}

		taskDto := models.CreateTaskDto{Type: taskType, Content: cryptorand.Text()}
		taskDtos = append(taskDtos, taskDto)
	}

	var wg sync.WaitGroup

	wg.Add(15)

	for _, taskDto := range taskDtos {
		go func() {
			randomDuration := rand.N(10 * time.Second)
			time.Sleep(randomDuration)
			defer wg.Done()

			postBody, _ := json.Marshal(taskDto)
			response, error := http.Post(
				"http://localhost:8000/queue",
				"application/json",
				bytes.NewBuffer(postBody),
			)

			if error != nil {
				log.Fatalf("Request with content %s failed", taskDto.Content)
			}

			defer response.Body.Close()

			_, err := io.ReadAll(response.Body)
			if err != nil {
				log.Fatalf("Error reading response for request with content %s", err)
			}
		}()
	}

	wg.Wait()

	response := map[string]string{
		"message": "Batch requests made to the Queue API",
	}

	json.NewEncoder(w).Encode(response)
}
