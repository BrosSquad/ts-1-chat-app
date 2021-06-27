package auth

import (
	"context"
)

type (
	LogoutService interface {
		Logout(context.Context, string) error
	}

	logoutService struct {
		service TokenService
	}
)

func (l logoutService) Logout(ctx context.Context, token string) error {
	return l.service.Delete(ctx, token)
}

func NewLogoutService(service TokenService) LogoutService {
	return logoutService{
		service: service,
	}
}
