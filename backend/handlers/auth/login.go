package auth

import (
	"context"
	
	"github.com/BrosSquad/ts-1-chat-app/backend/services/pb"
)

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
