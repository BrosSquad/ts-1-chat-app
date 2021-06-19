package services

import (
	"google.golang.org/grpc"
	"gorm.io/gorm"

	"github.com/BrosSquad/ts-1-chat-app/backend/services/auth"
	"github.com/BrosSquad/ts-1-chat-app/backend/services/chat"
	"github.com/BrosSquad/ts-1-chat-app/backend/services/pb"
)

func Register(server grpc.ServiceRegistrar, db *gorm.DB) {
	pb.RegisterChatServer(server, chat.New(db))
	pb.RegisterAuthServer(server, auth.New(db))
}
