package utils

import (
	"os"
	"strings"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	"github.com/gin-gonic/gin"
)

func GetEnv(key string, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

func GenerateResponse(httpStatusCode int, c *gin.Context, title1 string, message1 string, title2 string, input interface{}) {

	errorMessage := strings.TrimPrefix(message1, "ERROR: ")
	response := gin.H{
		title1: errorMessage,
	}

	if title2 != "" && input != nil {
		response[title2] = input
	}

	c.JSON(httpStatusCode, response)
}

func CreateSourceObj(combinedConfig model.CombinedStorageConfig) model.Source {
	var source model.Source

	source.SourcesName = combinedConfig.Source.SourcesName
	source.StorageType = combinedConfig.Source.StorageType

	return source
}

func CreateDestinationObj(combinedConfig model.CombinedStorageConfig) model.Destination {
	var destination model.Destination

	destination.DestinationsName = combinedConfig.Destination.DestinationsName
	destination.StorageType = combinedConfig.Destination.StorageType

	return destination
}

func CreateConfigObj(combinedConfig model.CombinedStorageConfig) model.Config {
	var config model.Config

	config.ClientID = combinedConfig.Config.ClientID
	config.ClientSecret = combinedConfig.Config.ClientSecret
	config.FolderURL = combinedConfig.Config.FolderURL
	config.RefreshToken = combinedConfig.Config.RefreshToken
	config.TokenURI = combinedConfig.Config.TokenURI

	return config
}

func CreatePipelineObj(sourceID int, destinationID int) model.Pipeline {
	var Pipeline model.Pipeline

	Pipeline.SourcesID = sourceID
	Pipeline.DestinationsID = destinationID

	return Pipeline
}
