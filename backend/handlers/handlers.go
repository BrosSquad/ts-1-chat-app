package handlers

import (
	"google.golang.org/grpc"

	"github.com/BrosSquad/ts-1-chat-app/backend/di"
	"github.com/BrosSquad/ts-1-chat-app/backend/handlers/auth"
	"github.com/BrosSquad/ts-1-chat-app/backend/handlers/chat"
	"github.com/BrosSquad/ts-1-chat-app/backend/services/pb"
)

func Register(server grpc.ServiceRegistrar, container di.Container) {
	db := container.GetDatabase()
	errorLogger := container.GetErrorLogger()
	debugLogger := container.GetDebugLogger()

	validator := container.GetValidator()

	pb.RegisterChatServer(server, chat.New(db, validator, debugLogger, errorLogger, container.GetChatBuffer()))

	pb.RegisterAuthServer(server, auth.New(container.GetRegisterService(), container.GetLoginService(), validator))
}
