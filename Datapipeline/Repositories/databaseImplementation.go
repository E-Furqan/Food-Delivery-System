package database

import (
	"log"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
)

func (repo *Repository) InsertSourceConfiguration(source *model.Source, config *model.Config) error {

	tx := repo.DB.Begin()
	if tx.Error != nil {
		log.Printf("Error starting transaction: %v", tx.Error)
		return tx.Error
	}

	err := tx.Create(source)
	if err.Error != nil {
		log.Printf("Error saving source: %v", err.Error)
		return err.Error
	}

	config.SourcesID = source.SourcesID
	err = tx.Create(config)
	if err.Error != nil {
		log.Printf("Error saving configuration: %v", err.Error)
		return err.Error
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing transaction: %v", err)
		return err
	}

	return nil
}

func (repo *Repository) InsertDestinationConfiguration(destination *model.Destination, config *model.Config) error {

	tx := repo.DB.Begin()
	if tx.Error != nil {
		log.Printf("Error starting transaction: %v", tx.Error)
		return tx.Error
	}

	err := tx.Create(destination)
	if err.Error != nil {
		log.Printf("Error saving source: %v", err.Error)
		return err.Error
	}

	config.DestinationsID = destination.DestinationsID
	err = tx.Create(config)
	if err.Error != nil {
		log.Printf("Error saving configuration: %v", err.Error)
		return err.Error
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing transaction: %v", err)
		return err
	}

	return nil
}

func (repo *Repository) CreatePipeline(pipeline *model.Pipeline) error {
	err := repo.DB.Create(pipeline)
	if err.Error != nil {
		log.Printf("Error creating pipeline: %v", err.Error)
		return err.Error
	}

	return nil
}

func (repo *Repository) FetchPipelineDetails(pipelineID int) (model.Pipeline, error) {

	var pipeline model.Pipeline

	err := repo.DB.Table("pipelines").
		Select("pipelines_id, sources_id , destinations_id").
		Where("pipelines_id = ?", pipeline).
		First(&pipeline).Error

	if err != nil {
		log.Printf("Error saving pipeline: %v", err)
		return model.Pipeline{}, err
	}

	return pipeline, nil
}

func (repo *Repository) FetchConfigSourceDetails(sourceID int) (model.Config, error) {

	var config model.Config

	err := repo.DB.Table("configs").
		Select("client_id, client_secret , token_uri,refresh_token,folder_Url").
		Where("sources_id = ?", sourceID).
		First(&config).Error

	if err != nil {
		log.Printf("Error saving pipeline: %v", err)
		return model.Config{}, err
	}

	return config, nil
}

func (repo *Repository) FetchConfigDestinationDetails(destinationID int) (model.Config, error) {

	var config model.Config

	err := repo.DB.Table("configs").
		Select("client_id, client_secret , token_uri,refresh_token,folder_Url").
		Where("destinations_id = ?", destinationID).
		First(&config).Error

	if err != nil {
		log.Printf("Error saving pipeline: %v", err)
		return model.Config{}, err
	}

	return config, nil
}

func (repo *Repository) AddLogs(logs model.Log) error {
	err := repo.DB.Create(logs)
	if err.Error != nil {
		log.Printf("Error creating logs: %v", err.Error)
		return err.Error
	}

	return nil
}
