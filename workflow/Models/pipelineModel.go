package model

type Pipeline struct {
	PipelinesID    int `json:"pipelines_id"`
	SourcesID      int `json:"sources_id"`
	DestinationsID int `json:"destinations_id"`
}
