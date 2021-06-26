package auth

import (
	"context"

	"github.com/BrosSquad/ts-1-chat-app/backend/models"
	"github.com/BrosSquad/ts-1-chat-app/backend/repositories/user"
	"github.com/BrosSquad/ts-1-chat-app/backend/services/password"
	"github.com/BrosSquad/ts-1-chat-app/backend/services/pb"
)

type (
	RegisterService interface {
		Register(ctx context.Context, req *pb.RegisterRequest) (models.User, error)
	}

	registerService struct {
		userRepository user.Repository
		hasher         password.Hasher
	}
)

func (r registerService) Register(ctx context.Context, req *pb.RegisterRequest) (user models.User, err error) {
	name := req.GetName()
	surname := req.GetSurname()
	email := req.GetEmail()
	pass := req.GetPassword()

	if err := r.userRepository.ExistsByEmail(ctx, email); err != nil {
		return user, err
	}

	user.Email = email
	user.Name = name
	user.Surname = surname
	user.Password = r.hasher.Hash(pass)

	user, err = r.userRepository.Create(ctx, user)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func NewRegisterService(userRepository user.Repository, hasher password.Hasher) RegisterService {
	return registerService{
		userRepository: userRepository,
		hasher:         hasher,
	}
}
