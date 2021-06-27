package auth

import (
	"context"

	"github.com/BrosSquad/ts-1-chat-app/backend/services/pb"
)

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
