package model

type Configs struct {
	ConfigID            int    `json:"configId"`
	ClientID            string `json:"clientId"`
	ClientSecret        string `json:"clientSecret"`
	TokenUri            string `json:"tokenUri"`
	AuthUri             string `json:"authUri"`
	RedirectUris        string `json:"redirectUris"`
	AuthProviderCertUrl string `json:"authProviderCertUrl"`
}
