package middleware

import (
	"context"
	"errors"
	"github.com/BrosSquad/ts-1-chat-app/backend/logging"
	"github.com/BrosSquad/ts-1-chat-app/backend/services/password"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/BrosSquad/ts-1-chat-app/backend/repositories"
	"github.com/BrosSquad/ts-1-chat-app/backend/services/pb"
	"github.com/BrosSquad/ts-1-chat-app/backend/validators"
)

func handleValidationError(validations validators.ValidationError) (err error) {
	st := status.New(codes.InvalidArgument, "invalid data")
	if validations.ValidationErrors != nil {
		l := len(validations.ValidationErrors)

		fields := make([]*pb.BadRequest_FieldViolation, 0, l)

		for key, validation := range validations.ValidationErrors {
			fields = append(fields, &pb.BadRequest_FieldViolation{
				Field:       key,
				Description: validation,
			})
		}

		badRequest := &pb.BadRequest{
			FieldViolations: fields,
		}

		st, err = st.WithDetails(badRequest)

		if err != nil {
			return status.Error(codes.Internal, "internal error")
		}

		return st.Err()
	}

	return status.Error(codes.InvalidArgument, "invalid data")
}

func handleRepositoryErrors(err error) error {
	switch {
	case errors.Is(err, repositories.ErrAlreadyExists):
		return status.Error(codes.AlreadyExists, "already exists")
	case errors.Is(err, repositories.ErrNotFound):
		return status.Error(codes.NotFound, "not found")
	case errors.Is(err, repositories.ErrDeleteFailed):
		return status.Error(codes.Internal, "delete failed")
	}

	if _, ok := err.(repositories.DatabaseError); ok {
		return status.Error(codes.Internal, "internal error")
	}

	return nil
}

func handleAuthErrors(err error) error {
	if errors.Is(err, password.ErrMismatchedHashAndPassword) {
		return status.Error(codes.Unauthenticated, "authentication error")
	}

	return nil
}

func handleError(err error) error {
	handled := handleRepositoryErrors(err)

	if handled != nil {
		return handled
	}

	handled = handleAuthErrors(err)

	if handled != nil {
		return handled
	}

	if validations, ok := err.(validators.ValidationError); ok {
		return handleValidationError(validations)
	}

	return status.Error(codes.Internal, "internal error")
}

func UnaryErrorHandler(errorLogger *logging.Error) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		res, err := handler(ctx, req)

		if err != nil {
			errorLogger.Err(err).
				Interface("request", req).
				Str("method", info.FullMethod).
				Msg("Something went wrong")
			return nil, handleError(err)
		}

		return res, nil
	}
}

func StreamErrorHandler() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		err := handler(srv, ss)

		if err != nil {
			return handleError(err)
		}

		return nil
	}
}
