package models

import "github.com/google/uuid"

type CreateTaskDto struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}
type Task struct {
	Uuid    uuid.UUID
	Type    int
	Content string
}
