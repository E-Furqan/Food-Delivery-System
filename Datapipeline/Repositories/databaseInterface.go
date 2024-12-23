package database

import (
	model "github.com/E-Furqan/Food-Delivery-System/Models"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

type RepositoryInterface interface {
	InsertSourceConfiguration(source *model.Source, config *model.Config) error
	InsertDestinationConfiguration(destination *model.Destination, config *model.Config) error
	CreatePipeline(pipeline model.Pipeline) error
	FetchPipelineDetails(pipelineID int) (model.Pipeline, error)
	FetchConfigSourceDetails(sourceID int) (model.Config, error)
	FetchConfigDestinationDetails(destinationID int) (model.Config, error)
	AddLogs(logs model.Log) error
}
