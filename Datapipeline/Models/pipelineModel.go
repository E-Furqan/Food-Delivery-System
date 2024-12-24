package model

type Pipeline struct {
	PipelineID     int `gorm:"primaryKey;autoIncrement" json:"pipeline_id"`
	SourcesID      int `json:"sources_id"`
	DestinationsID int `json:"destinations_id"`
}
