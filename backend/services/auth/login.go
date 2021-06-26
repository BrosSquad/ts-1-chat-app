package auth

import (
	"context"
	"errors"
	"github.com/BrosSquad/ts-1-chat-app/backend/logging"
	"github.com/BrosSquad/ts-1-chat-app/backend/models"
	"github.com/BrosSquad/ts-1-chat-app/backend/repositories/user"
	"github.com/BrosSquad/ts-1-chat-app/backend/services/password"
	"github.com/BrosSquad/ts-1-chat-app/backend/services/pb"
)

type (
	LoginService interface {
		Login(ctx context.Context, req *pb.LoginRequest) (models.User, string, error)
	}

	loginService struct {
		repo         user.Repository
		tokenService TokenService
		hasher       password.Hasher
		errorLogger  *logging.Error
	}
)

func NewLoginService(userRepo user.Repository, service TokenService, hasher password.Hasher, errorLogger *logging.Error) LoginService {
	return loginService{
		repo:         userRepo,
		tokenService: service,
		hasher:       hasher,
		errorLogger:  errorLogger,
	}
}

func (l loginService) Login(ctx context.Context, req *pb.LoginRequest) (models.User, string, error) {
	email, pass := req.GetEmail(), req.GetPassword()

	u, err := l.repo.FindByEmail(ctx, email)

	if err != nil {
		return models.User{}, "", err
	}

	if err := l.hasher.Verify(u.Password, pass); err != nil {
		if !errors.Is(err, password.ErrMismatchedHashAndPassword) {
			l.errorLogger.
				Err(err).
				Str("type", "password").
				Str("email", email).
				Msg("error while verifying password")
		}

		return models.User{}, "", password.ErrMismatchedHashAndPassword
	}

	token, err := l.tokenService.Generate(ctx, u)

	if err != nil {
		return models.User{}, "", err
	}

	return u, token, nil
}
