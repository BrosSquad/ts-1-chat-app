package auth

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/BrosSquad/ts-1-chat-app/backend/logging"
	"github.com/BrosSquad/ts-1-chat-app/backend/models"
	"github.com/BrosSquad/ts-1-chat-app/backend/services/password"
	"github.com/BrosSquad/ts-1-chat-app/backend/services/pb"
	"github.com/BrosSquad/ts-1-chat-app/backend/validators"
)

type authService struct {
	errorLogger    *logging.Error
	db             *gorm.DB
	passwordHasher password.Hasher
	validator      validators.Validator

	pb.UnimplementedAuthServer
}

func New(db *gorm.DB, errorLogger *logging.Error, hasher password.Hasher, validator validators.Validator) pb.AuthServer {
	return &authService{
		errorLogger:    errorLogger,
		db:             db,
		passwordHasher: hasher,
		validator:      validator,
	}
}

func (a *authService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	validations, err := a.validator.Struct(req)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid input")
	}

	if validations != nil {
		return &pb.RegisterResponse{
			ValidationErrors: validations,
			User:             nil,
		}, status.Error(codes.InvalidArgument, "invalid input")
	}

	name := req.GetName()
	surname := req.GetSurname()
	email := req.GetEmail()
	pass := req.GetPassword()

	var user models.User

	tx := a.db.WithContext(ctx)

	result := tx.Where("email = ?", email).First(&user)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			a.errorLogger.
				Err(result.Error).
				Str("type", "database").
				Str("query", result.Statement.SQL.String()).
				Msg("error while reading the database")

			return nil, status.Error(codes.Internal, "cannot insert user")
		}

		user.Email = email
		user.Name = name
		user.Surname = surname
		user.Password = a.passwordHasher.Hash(pass)

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
			Id:      user.ID,
			Email:   user.Email,
			Name:    name,
			Surname: surname,
		},
	}, nil
}

func (a *authService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	validations, err := a.validator.Struct(req)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid input")
	}

	if validations != nil {
		return &pb.LoginResponse{
			ValidationErrors: validations,
		}, status.Error(codes.InvalidArgument, "invalid input")
	}

	email, pass := req.GetEmail(), req.GetPassword()

	var user models.User

	tx := a.db.WithContext(ctx)

	result := tx.Where("email = ?", email).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "User is not found")
		}

		a.errorLogger.
			Err(result.Error).
			Str("type", "database").
			Str("email", email).
			Str("query", result.Statement.SQL.String()).
			Msg("error while reading the database")

		return nil, status.Error(codes.Internal, "cannot fetch user")
	}

	err = a.passwordHasher.Verify(user.Password, pass)

	if err != nil {
		if errors.Is(err, password.ErrMismatchedHashAndPassword) {
			a.errorLogger.
				Err(err).
				Str("type", "password").
				Str("email", email).
				Msg("error while verifiein password")

			return nil, status.Error(codes.Internal, "password error")
		}

		return nil, status.Error(codes.Unauthenticated, "mismatched credentials")
	}

	// TODO: Generate Token

	return &pb.LoginResponse{
		User: & pb.User{
			Id:      user.ID,
			Name:    user.Name,
			Surname: user.Surname,
			Email:   user.Email,
		},
		Token:            "",
		ValidationErrors: nil,
	}, nil
}
