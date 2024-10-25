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
	Username    string `json:"username"`
	ActiveRole  string `json:"activeRole"`
	ServiceType string `json:"service_type"`
	jwt.StandardClaims
}

type Claims interface {
	jwt.Claims
	SetExpirationTime(expiration int64)
}

type UserClaims struct {
	Username    string `json:"username"`
	ActiveRole  string `json:"activeRole"`
	ServiceType string `json:"service_type"`
	jwt.StandardClaims
}

func (u *UserClaims) SetExpirationTime(expiration int64) {
	u.ExpiresAt = expiration
}

type RestaurantClaims struct {
	ClaimId     uint   `json:"claim_id"`
	ServiceType string `json:"service_type"`
	jwt.StandardClaims
}

func (r *RestaurantClaims) SetExpirationTime(expiration int64) {
	r.ExpiresAt = expiration
}

type OrderClaims struct {
	ClaimId     uint   `json:"claim_id"`
	ServiceType string `json:"service_type"`
	jwt.StandardClaims
}

func (o *OrderClaims) SetExpirationTime(expiration int64) {
	o.ExpiresAt = expiration
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expiration   int64  `json:"expires_at"`
}
