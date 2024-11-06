package model

import "github.com/golang-jwt/jwt"

type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
	Role         string `json:"activeRole"`
}

type RestaurantClaim struct {
	ClaimId uint   `json:"claim_id"`
	Role    string `json:"activeRole"`
	jwt.StandardClaims
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expiration   int64  `json:"expires_at"`
}
