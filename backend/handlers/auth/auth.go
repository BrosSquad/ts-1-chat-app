package auth

import (
	"github.com/BrosSquad/ts-1-chat-app/backend/services/auth"
	"github.com/BrosSquad/ts-1-chat-app/backend/services/pb"
	"github.com/BrosSquad/ts-1-chat-app/backend/validators"
)

type authService struct {
	registerService auth.RegisterService
	loginService    auth.LoginService
	logoutService   auth.LogoutService
	validator       validators.Validator

	pb.UnimplementedAuthServer
}

func New(registerService auth.RegisterService, logoutService auth.LogoutService, loginService auth.LoginService, validator validators.Validator) pb.AuthServer {
	return &authService{
		registerService: registerService,
		loginService:    loginService,
		logoutService:    logoutService,
		validator:       validator,
	}
}
