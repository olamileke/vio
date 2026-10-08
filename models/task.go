package models

import "github.com/google/uuid"

type Task struct {
	Uuid    uuid.UUID
	Type    int
	Content string
}
