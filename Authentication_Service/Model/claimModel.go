package model

import "github.com/golang-jwt/jwt"

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
