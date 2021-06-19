package chat

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/BrosSquad/ts-1-chat-app/backend/models"
	"github.com/BrosSquad/ts-1-chat-app/backend/services/pb"
)

type UserMessage struct {
	User    models.User
	Message models.Message
}

type chatService struct {
	messages chan UserMessage
	db       *gorm.DB
	pb.UnimplementedChatServer
}

func New(db *gorm.DB) pb.ChatServer {
	return &chatService{
		db:       db,
		messages: make(chan UserMessage, 5000),
	}
}

func (c *chatService) SendMessage(ctx context.Context, in *pb.MessageRequest) (*pb.Empty, error) {
	message := models.Message{
		UserID: in.GetUserId(),
		Text:   in.GetText(),
	}

	result := c.db.Create(&message)

	if result.Error != nil {
		return nil, status.Error(codes.Internal, "cannot insert message")
	}

	var user models.User

	result = c.db.Where("id = ?", in.GetUserId()).First(&user)

	if result.Error != nil {
		return nil, status.Error(codes.Internal, "cannot fetch a user")
	}

	c.messages <- UserMessage{
		User:    user,
		Message: message,
	}

	return &pb.Empty{}, nil
}

func (c *chatService) Connect(req *pb.ConnectRequest, client pb.Chat_ConnectServer) error {
	messages := make([]models.Message, 0, 50)
	result := c.db.Order("created_at DESC").Limit(50).Find(&messages)

	if result.Error != nil {
		return status.Error(codes.Internal, "Cannot fetch messages")
	}

	for _, message := range messages {
		client.Send(&pb.MessageResponse{
			User: &pb.User{
				Id:       message.User.ID,
				Username: message.User.Username,
			},
			Text:      message.Text,
			CreatedAt: message.CreatedAt.Format(time.RFC3339),
		})
	}

	for {
		select {
		case um := <-c.messages:
			client.Send(&pb.MessageResponse{
				User: &pb.User{
					Id:       um.User.ID,
					Username: um.User.Username,
				},
				Text:      um.Message.Text,
				CreatedAt: um.Message.CreatedAt.Format(time.RFC3339),
			})

		case <-client.Context().Done():
			return nil
		}
	}
}
