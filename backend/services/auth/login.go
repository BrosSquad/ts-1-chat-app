package auth

import (
	"github.com/BrosSquad/ts-1-chat-app/backend/repositories/user"
	"github.com/BrosSquad/ts-1-chat-app/backend/services/password"
)

type (
	LoginService interface{}

	loginService struct{}
)

func NewLoginService(userRepo user.Repository, tokenRepo TokenService, hasher password.Hasher) LoginService {
	return loginService{}
}
