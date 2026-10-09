package models

import "github.com/google/uuid"

type Result struct {
	Uuid     uuid.UUID
	Type     int
	Content  string
	Status   string
	Duration float64
}
