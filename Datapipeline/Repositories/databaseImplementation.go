package database

import (
	"log"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
)

func (repo *Repository) InsertSourceConfiguration(combinedConfig *model.CombinedStorageConfig) error {
	source := utils.CreateSourceObj(combinedConfig)

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

	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing transaction: %v", err)
		return err
	}

	config := utils.CreateConfigObj(combinedConfig)
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

func (repo *Repository) InsertDestinationConfiguration(combinedConfig *model.CombinedStorageConfig) error {
	destination := utils.CreateDestinationObj(combinedConfig)

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

	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing transaction: %v", err)
		return err
	}

	config := utils.CreateConfigObj(combinedConfig)
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

func (repo *Repository) InsertPipeline(sourceID int, destinationID int) error {
	pipeline := utils.CreatePipelineObj(sourceID, destinationID)

	err := repo.DB.Create(pipeline)
	if err.Error != nil {
		log.Printf("Error creating pipeline: %v", err.Error)
		return err.Error
	}

	return nil
}
