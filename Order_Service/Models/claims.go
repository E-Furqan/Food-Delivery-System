package model

import "github.com/golang-jwt/jwt"

type Claims struct {
	UserId      uint   `json:"user_id"`
	ActiveRole  string `json:"activeRole"`
	ClaimId     uint   `json:"claim_id"`
	ServiceType string `json:"service_type"`
	jwt.StandardClaims
}
