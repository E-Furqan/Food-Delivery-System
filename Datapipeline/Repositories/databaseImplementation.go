package database

import (
	"log"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
)

func (repo *Repository) InsertConfiguration(configuration *model.Configs) error {
	result := repo.DB.Create(configuration)
	if result.Error != nil {
		log.Printf("Error saving configuration: %v", result.Error)
		return result.Error
	}
	return nil
}
