package model

import "time"

type Token struct {
	TokenID      uint      `json:"token_id"`
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	RefreshToken string    `json:"refresh_token"`
	Expiry       time.Time `json:"expiry"`
}

type TokenConfig struct {
	TokenID  uint `json:"token_id"`
	ConfigID uint `json:"config_id"`
}
