package chat

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/BrosSquad/ts-1-chat-app/backend/models"
	"github.com/BrosSquad/ts-1-chat-app/backend/services/pb"
)

type chatService struct {
	db *gorm.DB
	pb.UnimplementedChatServer
}

func New(db *gorm.DB) pb.ChatServer {
	return &chatService{
		db: db,
	}
}

func (c *chatService) SendMessage(ctx context.Context, in *pb.MessageRequest) (*pb.Empty, error) {
	message := models.Message{
		UserID: in.GetUserId(),
		Text: in.GetText(),
	}

	result := c.db.Create(&message)

	if result.Error != nil {
		return nil, status.Error(codes.Internal, "cannot insert user")
	}

	return &pb.Empty{}, nil
}

func (c *chatService) Connect(req *pb.ConnectRequest, client pb.Chat_ConnectServer) error {
	return nil
}
