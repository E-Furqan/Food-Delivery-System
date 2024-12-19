package model

import "time"

type Token struct {
	TokenID      uint
	AccessToken  string
	TokenType    string
	RefreshToken string
	Expiry       time.Time
	CreatedAt    time.Time
}
