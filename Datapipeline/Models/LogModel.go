package model

import "time"

type Log struct {
	LogID      int       `json:"logId"`
	LogMessage string    `json:"logMessage"`
	CreatedAt  time.Time `json:"createdAt"`
}

type LogConfig struct {
	LogID    int `json:"logId"`
	ConfigID int `json:"configId"`
}
