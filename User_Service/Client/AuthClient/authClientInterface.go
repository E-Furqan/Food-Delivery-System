package AuthClient

import model "github.com/E-Furqan/Food-Delivery-System/Models"

type AuthClient struct {
	model.AuthClientEnv
}

func NewClient(env model.AuthClientEnv) *AuthClient {
	return &AuthClient{
		AuthClientEnv: env,
	}
}

type AuthClientInterface interface {
	GenerateToken(input model.UserClaim) (*model.Tokens, error)
	RefreshToken(input model.RefreshToken) (*model.Tokens, error)
}
