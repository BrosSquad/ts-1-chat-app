package chat

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/BrosSquad/ts-1-chat-app/backend/models"
	"github.com/BrosSquad/ts-1-chat-app/backend/services/pb"
)

func (c *chatService) SendMessage(ctx context.Context, in *pb.MessageRequest) (*pb.MessageResponse, error) {
	err := c.validator.Struct(in)

	if err != nil {
		return nil, err
	}

	message := models.Message{
		UserID: in.GetUserId(),
		Text:   in.GetText(),
	}

	result := c.db.
		WithContext(ctx).
		Create(&message)

	if result.Error != nil {
		c.errorLogger.
			Err(result.Error).
			Str("type", "database").
			Str("query", result.Statement.SQL.String()).
			Str("table", result.Statement.Table).
			Str("message", message.Text).
			Uint64("userId", message.UserID).
			Msg("cannot insert message into database")

		return nil, status.Error(codes.Internal, "cannot insert message")
	}

	var user models.User

	result = c.db.Where("id = ?", in.GetUserId()).First(&user)

	if result.Error != nil {
		c.errorLogger.
			Err(result.Error).
			Str("type", "database").
			Str("query", result.Statement.SQL.String()).
			Str("table", result.Statement.Table).
			Uint64("userId", message.UserID).
			Msg("cannot fetch user from database")
		return nil, status.Error(codes.Internal, "cannot fetch a user")
	}

	c.messages <- UserMessage{
		User:    user,
		Message: message,
	}

	return &pb.MessageResponse{
		User: &pb.User{
			Id:      user.ID,
			Name:    user.Name,
			Surname: user.Surname,
			Email:   user.Email,
		},
		Text:      message.Text,
		CreatedAt: message.CreatedAt.Format(time.RFC3339),
	}, nil
}
