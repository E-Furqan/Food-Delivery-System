package model

type Config struct {
	ClientID       string `json:"client_id"`
	ClientSecret   string `json:"client_secret"`
	TokenURI       string `json:"token_uri"`
	RefreshToken   string `json:"refresh_token"`
	SourcesID      int    `json:"sources_id"`
	DestinationsID int    `json:"destinations_id"`
	FolderURL      string `json:"folder_url"`

	Source      Source      `json:"source"`
	Destination Destination `json:"destination"`
}

type Source struct {
	SourcesID   int    `json:"sources_id"`
	SourcesName string `json:"sources_name"`
	SourceType  string `json:"source_type"`
}

type Destination struct {
	DestinationsID   int    `json:"destinations_id"`
	DestinationsName string `json:"destinations_name"`
	DestinationType  string `json:"destination_type"`
}
