package model

type Configs struct {
	ConfigID            uint   `json:"config_id"`
	ClientID            string `json:"client_id"`
	ClientSecret        string `json:"client_secret"`
	TokenUri            string `json:"token_uri"`
	AuthUri             string `json:"auth_uri"`
	RedirectUris        string `json:"redirect_uris"`
	AuthProviderCertUrl string `json:"auth_provider_cert_url"`
}
