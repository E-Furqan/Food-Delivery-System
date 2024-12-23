package model

type Source struct {
	SourcesID   int    `json:"sources_id"`
	SourcesName string `json:"sources_name"`
	StorageType string `json:"storage_type" binding:"required"`
}

type Destination struct {
	DestinationsID   int    `json:"destinations_id"`
	DestinationsName string `json:"destinations_name"`
	StorageType      string `json:"storage_type" binding:"required"`
}

type Config struct {
	ClientID       string `json:"client_id" binding:"required"`
	ClientSecret   string `json:"client_secret" binding:"required"`
	TokenURI       string `json:"token_uri" binding:"required"`
	RefreshToken   string `json:"refresh_token" binding:"required"`
	FolderURL      string `json:"folder_url" binding:"required"`
	SourcesID      int    `json:"sources_id"`
	DestinationsID int    `json:"destinations_id"`
}
