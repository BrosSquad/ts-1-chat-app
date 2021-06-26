package auth

import "github.com/BrosSquad/ts-1-chat-app/backend/models"

type (
	TokenService interface {
		Generate() models.Token
		Verify(string) error
	}

	tokenService struct{}
)

func (t tokenService) Generate() models.Token {
	return models.Token{}
}

func (t tokenService) Verify(token string) error {
	return nil
}

func NewTokenService() TokenService {
	return &tokenService{}
}
