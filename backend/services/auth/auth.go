package auth

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/BrosSquad/ts-1-chat-app/backend/models"
	"github.com/BrosSquad/ts-1-chat-app/backend/services/pb"
)

type authService struct{
	db *gorm.DB
	pb.UnimplementedAuthServer
}

func New(db *gorm.DB) pb.AuthServer {
	return &authService{
		db: db,
	}
}

func (a *authService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	user := models.User {
		Username: req.GetUsername(),
	}

	result := a.db.Create(&user)

	if result.Error != nil {
		return nil, status.Error(codes.Internal, "cannot insert user")
	}

	return &pb.RegisterResponse{
		User: &pb.User{
			Id: user.ID,
			Username: user.Username,
		},
	}, nil
}
