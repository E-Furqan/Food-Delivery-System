package payload

import (
	"github.com/golang-jwt/jwt"
)

type Environment struct {
	JWT_SECRET      string
	RefreshTokenKey string
}

type Input struct {
	ClaimId     uint   `json:"claim_id"`
	UserId      uint   `json:"user_id"`
	ActiveRole  string `json:"activeRole"`
	ServiceType string `json:"service_type"`
	jwt.StandardClaims
}

type Claims interface {
	jwt.Claims
	SetExpirationTime(expiration int64)
}

type UserClaims struct {
	UserId      uint   `json:"user_id"`
	ActiveRole  string `json:"activeRole"`
	ServiceType string `json:"service_type"`
	jwt.StandardClaims
}

func (u *UserClaims) SetExpirationTime(expiration int64) {
	u.ExpiresAt = expiration
}

type IDClaims struct {
	ClaimId     uint   `json:"claim_id"`
	ServiceType string `json:"service_type"`
	jwt.StandardClaims
}

func (o *IDClaims) SetExpirationTime(expiration int64) {
	o.ExpiresAt = expiration
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expiration   int64  `json:"expires_at"`
}
