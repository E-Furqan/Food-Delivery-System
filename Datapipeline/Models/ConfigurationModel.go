package model

type CombinedSourceStorageConfig struct {
	Config `json:"config"`
	Source `json:"source"`
}

type CombinedDestinationStorageConfig struct {
	Config      `json:"config"`
	Destination `json:"destination"`
}

type Source struct {
	SourcesID   int    `gorm:"primaryKey;autoIncrement" json:"sources_id"`
	SourcesName string `json:"sources_name"`
	StorageType string `json:"storage_type"`
}

type Destination struct {
	DestinationsID   int    `gorm:"primaryKey;autoIncrement" json:"destinations_id"`
	DestinationsName string `json:"destinations_name"`
	StorageType      string `json:"storage_type"`
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
