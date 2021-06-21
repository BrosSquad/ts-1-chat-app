package auth

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/BrosSquad/ts-1-chat-app/backend/logging"
	"github.com/BrosSquad/ts-1-chat-app/backend/models"
	"github.com/BrosSquad/ts-1-chat-app/backend/services/pb"
)

type authService struct {
	errorLogger *logging.Error
	db          *gorm.DB
	pb.UnimplementedAuthServer
}

func New(db *gorm.DB, errorLogger *logging.Error) pb.AuthServer {
	return &authService{
		errorLogger: errorLogger,
		db:          db,
	}
}

func (a *authService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	username := req.GetUsername()

	var user models.User

	tx := a.db.WithContext(ctx)

	result := tx.Where("username = ?", username).First(&user)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			a.errorLogger.
				Err(result.Error).
				Str("type", "database").
				Str("query", result.Statement.SQL.String()).
				Msg("error while reading the database")

			return nil, status.Error(codes.Internal, "cannot insert user")
		}

		user.Username = username

		result = tx.Create(&user)

		if result.Error != nil {
			a.errorLogger.
				Err(result.Error).
				Str("type", "database").
				Str("query", result.Statement.SQL.String()).
				Msg("cannot create new user")

			return nil, status.Error(codes.Internal, "cannot insert user")
		}
	}

	return &pb.RegisterResponse{
		User: &pb.User{
			Id:       user.ID,
			Username: user.Username,
		},
	}, nil
}
