package model

import "time"

type Log struct {
	LogID       int       `gorm:"primaryKey;autoIncrement;column:log_id" json:"log_id"`
	LogMessage  string    `gorm:"type:text;column:log_message" json:"log_message"`
	PipelinesID int       `gorm:"column:pipelines_id" json:"pipelines_id"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP;column:created_at" json:"created_at"`
	Pipeline    Pipeline  `json:"pipeline"`
}

type LogConfig struct {
	LogID    uint `json:"log_id"`
	ConfigID uint `json:"config_id"`
}
