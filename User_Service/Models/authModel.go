package model

import "github.com/golang-jwt/jwt"

type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
	Role         string `json:"activeRole"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expiration   int64  `json:"expires_at"`
}

type Claims struct {
	ClaimId uint   `json:"claim_id"`
	Role    string `json:"activeRole"`
	jwt.StandardClaims
}

type MiddlewareEnv struct {
	JWT_SECRET      string
	RefreshTokenKey string
}

type AuthClientEnv struct {
	BASE_URL           string
	GENERATE_TOKEN_URL string
	REFRESH_TOKEN_URL  string
	AUTH_PORT          string
}
