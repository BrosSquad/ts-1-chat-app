package middleware

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"github.com/BrosSquad/ts-1-chat-app/backend/logging"
)

func UnaryElapsed(debugLogger *logging.Debug) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := time.Now().UTC()

		resp, err = handler(ctx, req)

		dur := time.Now().UTC().Sub(start)

		debugLogger.Debug().Str("elapsed", dur.String()).Msg("Request time")

		return resp, err
	}
}
