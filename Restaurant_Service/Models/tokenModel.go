package model

type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
	ServiceType  string `json:"service_type"`
}
type RestaurantClaim struct {
	ClaimId     uint   `json:"claim_id"`
	ServiceType string `json:"service_type"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expiration   int64  `json:"expires_at"`
}
