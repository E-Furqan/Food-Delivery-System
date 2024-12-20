package database

import (
	"fmt"
	"log"
	"time"

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

func (repo *Repository) DeleteExpiredTokens() {
	result := repo.DB.Where("expiry <= ?", time.Now()).Delete(&model.Token{})
	if result.Error != nil {
		fmt.Println("Error deleting expired tokens:", result.Error)
	} else {
		fmt.Printf("Deleted %d expired tokens.\n", result.RowsAffected)
	}
}

func (repo *Repository) InsertToken(token *model.Token, configID uint) error {
	result := repo.DB.Create(token)
	if result.Error != nil {
		log.Printf("Error saving token: %v", result.Error)
		return result.Error
	}
	repo.LinkTokenToConfig(token.TokenID, configID)
	return nil
}

func (repo *Repository) LinkTokenToConfig(tokenID uint, configID uint) error {
	tokenConfig := model.TokenConfig{
		TokenID:  tokenID,
		ConfigID: configID,
	}

	result := repo.DB.Create(&tokenConfig)
	if result.Error != nil {
		log.Printf("Error linking token to configuration: %v", result.Error)
		return result.Error
	}

	return nil
}

func (repo *Repository) FetchToken(token *model.Token) error {
	result := repo.DB.Table("tokens").
		Select("access_token, token_type, refresh_token, expiry").
		Where("token_id = ?", token.TokenID).
		First(token)

	if result.Error != nil {
		log.Printf("Error saving token: %v", result.Error)
		return result.Error
	}

	log.Print("token from database ", token)

	return nil
}

func (repo *Repository) LinkLogsToConfig(logID uint, configID uint) error {
	logConfig := model.LogConfig{
		LogID:    logID,
		ConfigID: configID,
	}

	result := repo.DB.Create(&logConfig)
	if result.Error != nil {
		log.Printf("Error linking token to configuration: %v", result.Error)
		return result.Error
	}

	return nil
}
