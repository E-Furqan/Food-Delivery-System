package model

import "time"

type Log struct {
	LogID       int       `json:"log_id"`
	LogMessage  string    `json:"log_message"`
	PipelinesID int       `json:"pipelines_id"`
	CreatedAt   time.Time `json:"created_at"`
}
