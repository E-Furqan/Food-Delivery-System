package model

type Configs struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	TokenUri     string `json:"token_uri"`
	RefreshToken string `json:"refresh_token"`
}
