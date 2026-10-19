package controllers

import (
	"encoding/json"
	"net/http"
	"strings"
	"vio/models"
	"vio/queue"

	"github.com/google/uuid"
)

type queueController struct {
	queue *queue.Queue
}

func RegisterQueueRoutes(mux *http.ServeMux, queue *queue.Queue) {
	qc := queueController{queue: queue}
	mux.HandleFunc("/queue", qc.enqueue)
}

func (qc *queueController) enqueue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var createTaskDto models.CreateTaskDto

	if err := json.NewDecoder(r.Body).Decode(&createTaskDto); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	var parsedType int

	switch strings.ToLower(createTaskDto.Type) {
	case "email":
		parsedType = 1
	case "pdf":
		parsedType = 2
	default:
		parsedType = 1
	}

	task := models.Task{Uuid: uuid.New(), Type: parsedType, Content: createTaskDto.Content}

	qc.queue.Process(task)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"message": "Task queued successfully"}
	json.NewEncoder(w).Encode(response)
}
