package middleware

import (
	"google.golang.org/grpc"
	"sort"

	"github.com/BrosSquad/ts-1-chat-app/backend/di"
)

func Register(c di.Container) []grpc.ServerOption {
	allowedList := c.GetConfig().GetStringSlice("auth.allowed")
	blockedList := c.GetConfig().GetStringSlice("auth.blocked")

	sort.Strings(allowedList)
	sort.Strings(blockedList)

	return []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			UnaryElapsed(c.GetDebugLogger()),
			UnaryErrorHandler(c.GetErrorLogger()),
			UnaryAuth(
				AuthConfig{
					Allowed: allowedList,
					Blocked: blockedList,
				},
				c.GetTokenService(),
			),
		),
		grpc.ChainStreamInterceptor(
			StreamErrorHandler(),
			StreamAuth(
				AuthConfig{
					Allowed: allowedList,
					Blocked: blockedList,
				},
				c.GetTokenService(),
			),
		),
	}
}
