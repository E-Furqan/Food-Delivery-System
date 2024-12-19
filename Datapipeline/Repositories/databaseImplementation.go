package database

import (
	"fmt"
	"log"
	"time"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
)

func (repo *Repository) CreateConfiguration(configuration *model.Configs) error {
	result := repo.DB.Create(configuration)
	if result.Error != nil {
		log.Printf("Error saving configuration: %v", result.Error)
		return result.Error
	}
	return nil
}

func (repo *Repository) DeleteExpiredTokens() {
	// Delete tokens where the expiry time has passed
	result := repo.DB.Where("expiry <= ?", time.Now()).Delete(&model.Token{})
	if result.Error != nil {
		fmt.Println("Error deleting expired tokens:", result.Error)
	} else {
		fmt.Printf("Deleted %d expired tokens.\n", result.RowsAffected)
	}
}
