package controllers

import (
	"encoding/json"
	"net/http"
	"vio/models"
	"vio/queue"

	"github.com/google/uuid"
)

type queueController struct{}

func RegisterQueueRoutes(mux *http.ServeMux) {
	qc := queueController{}
	mux.HandleFunc("/queue", qc.queue)
}

func (qc *queueController) queue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var createTaskDto models.CreateTaskDto

	if err := json.NewDecoder(r.Body).Decode(&createTaskDto); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	task := models.Task{Uuid: uuid.New(), Type: createTaskDto.Type, Content: createTaskDto.Content}

	queue.Queue(task)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"message": "Task queued successfully"}
	json.NewEncoder(w).Encode(response)
}
