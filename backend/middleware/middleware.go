package middleware

import (
	"google.golang.org/grpc"

	"github.com/BrosSquad/ts-1-chat-app/backend/di"
)

func Register(c di.Container) []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			UnaryElapsed(c.GetDebugLogger()),
			UnaryErrorHandler(c.GetErrorLogger()),
		),
		grpc.ChainStreamInterceptor(
			StreamErrorHandler(),
		),
	}
}
