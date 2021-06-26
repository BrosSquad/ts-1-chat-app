package middleware

import "google.golang.org/grpc"

func Register() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			UnaryErrorHandler(),
		),
		grpc.ChainStreamInterceptor(
			StreamErrorHandler(),
		),
	}
}
