package middleware

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sort"

	"github.com/BrosSquad/ts-1-chat-app/backend/services/auth"
)

type AuthConfig struct {
	Allowed sort.StringSlice
	Blocked sort.StringSlice
}

func authCheck(ctx context.Context, method string, config *AuthConfig, tokenService auth.TokenService) error {
	blockedLen := len(config.Blocked)
	allowedLen := len(config.Allowed)

	if blockedLen != 0 {
		idx := sort.Search(blockedLen, func(i int) bool {
			return config.Blocked[i] == method
		})

		// Found in blacklist
		if idx != blockedLen {
			tokenType, token, err := auth.ExtractToken(ctx)

			if err != nil {
				return err
			}

			return tokenService.Verify(ctx, tokenType, token)
		}
	}

	if allowedLen == 0 {
		// Whitelist is empty, everything should be allowed
		return nil

	}

	idx := sort.Search(allowedLen, func(i int) bool {
		return config.Allowed[i] == method
	})

	// Found in whitelist
	if idx != allowedLen {
		return nil
	}

	tokenType, token, err := auth.ExtractToken(ctx)

	if err != nil {
		return err
	}

	// If not found in whitelist, then token must be checked
	return tokenService.Verify(ctx, tokenType, token)
}

func shouldCheck(config AuthConfig) bool {
	blockedLen := len(config.Blocked)
	allowedLen := len(config.Allowed)

	shouldCheck := true

	// Both lists are empty, everything should go through
	if blockedLen == 0 && allowedLen == 0 {
		shouldCheck = false
	}

	return shouldCheck
}

func UnaryAuth(config AuthConfig, tokenService auth.TokenService) grpc.UnaryServerInterceptor {
	check := shouldCheck(config)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if !check {
			return handler(ctx, req)
		}

		method, exist := grpc.Method(ctx)

		if !exist {
			return nil, status.Error(codes.Unauthenticated, "unauthenticated")
		}

		if err := authCheck(ctx, method, &config, tokenService); err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		return handler(ctx, req)
	}
}

func StreamAuth(config AuthConfig, tokenService auth.TokenService) grpc.StreamServerInterceptor {
	check := shouldCheck(config)

	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		if !check {
			return handler(srv, ss)
		}

		method, exist := grpc.MethodFromServerStream(ss)

		if !exist {
			return status.Error(codes.Unauthenticated, "unauthenticated")
		}

		if err := authCheck(ss.Context(), method, &config, tokenService); err != nil {
			return status.Error(codes.Unauthenticated, "unauthenticated")
		}

		return handler(srv, ss)
	}
}
