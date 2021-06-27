package auth

import (
	"context"

	"github.com/BrosSquad/ts-1-chat-app/backend/services/auth"
	"github.com/BrosSquad/ts-1-chat-app/backend/services/pb"
)

func (a *authService) Logout(ctx context.Context, _ *pb.Empty) (*pb.Empty, error) {
	_, token, err := auth.ExtractToken(ctx)

	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, a.logoutService.Logout(ctx, token)
}
