package models

import "time"

type Log struct {
	Type      int       `json:"type"`
	Content   string    `json:"content"`
	Status    int       `json:"status"`
	Duration  float64   `json:"duration"`
	Timestamp time.Time `json:"timestamp"`
	Error     *string   `json:"error"`
}

type LogData struct {
	Logs []Log `json:"logs"`
}
