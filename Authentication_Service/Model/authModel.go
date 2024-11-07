package model

type AuthSecrets struct {
	JWT_SECRET      string
	RefreshTokenKey string
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expiration   int64  `json:"expires_at"`
}
