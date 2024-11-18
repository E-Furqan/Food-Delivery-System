package payload

import (
	"github.com/golang-jwt/jwt"
)

type Environment struct {
	JWT_SECRET      string
	RefreshTokenKey string
}

type Input struct {
	ClaimId    uint   `json:"claim_id"`
	ActiveRole string `json:"activeRole"`
	jwt.StandardClaims
}

type Claims interface {
	jwt.Claims
	SetExpirationTime(expiration int64)
}

type GeneralClaim struct {
	ClaimId    uint   `json:"claim_id"`
	ActiveRole string `json:"activeRole"`
	jwt.StandardClaims
}

func (u *GeneralClaim) SetExpirationTime(expiration int64) {
	u.ExpiresAt = expiration
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expiration   int64  `json:"expires_at"`
}
