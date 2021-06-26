package auth

import (
	"context"
	"github.com/BrosSquad/ts-1-chat-app/backend/services/pb"
	"github.com/BrosSquad/ts-1-chat-app/backend/validators"
)

type authService struct {
	registerService RegisterService
	loginService    LoginService
	validator       validators.Validator

	pb.UnimplementedAuthServer
}

func New(registerService RegisterService, loginService LoginService, validator validators.Validator) pb.AuthServer {
	return &authService{
		registerService: registerService,
		loginService:    loginService,
		validator:       validator,
	}
}

func (a *authService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	err := a.validator.Struct(req)

	if err != nil {
		return nil, err
	}

	user, err := a.registerService.Register(ctx, req)

	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{
		User: &pb.User{
			Id:      user.ID,
			Email:   user.Email,
			Name:    user.Name,
			Surname: user.Surname,
		},
	}, nil
}

func (a *authService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	err := a.validator.Struct(req)

	if err != nil {
		return nil, err
	}

	user, token, err := a.loginService.Login(ctx, req)

	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		User: &pb.User{
			Id:      user.ID,
			Name:    user.Name,
			Surname: user.Surname,
			Email:   user.Email,
		},
		Token: token,
	}, nil
}
