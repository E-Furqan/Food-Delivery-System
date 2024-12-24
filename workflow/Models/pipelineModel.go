package model

type Pipeline struct {
	PipelineID     int `json:"pipeline_id"`
	SourcesID      int `json:"sources_id"`
	DestinationsID int `json:"destinations_id"`
}
