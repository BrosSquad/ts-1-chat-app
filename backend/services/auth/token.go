package auth

import (
	"github.com/BrosSquad/ts-1-chat-app/backend/models"
	"github.com/BrosSquad/ts-1-chat-app/backend/repositories/token"
)

type (
	TokenService interface {
		Generate() models.Token
		Verify(string) error
	}

	tokenService struct {
		repo token.Repository
	}
)

func (t tokenService) Generate() models.Token {
	return models.Token{}
}

func (t tokenService) Verify(token string) error {
	return nil
}

func NewTokenService(repo token.Repository) TokenService {
	return &tokenService{
		repo: repo,
	}
}
