package model

type Configuration struct {
	ConfigID            int    `json:"configId"`
	ClientID            string `json:"clientId"`
	ClientSecret        string `json:"clientSecret"`
	TokenUri            string `json:"tokenUri"`
	AuthUri             string `json:"authUri"`
	RedirectUris        string `json:"redirectUris"`
	AuthProviderCertUrl string `json:"authProviderCertUrl"`
}
