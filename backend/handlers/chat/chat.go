package chat

import (
	"gorm.io/gorm"
	"sync"

	"github.com/BrosSquad/ts-1-chat-app/backend/logging"
	"github.com/BrosSquad/ts-1-chat-app/backend/models"
	"github.com/BrosSquad/ts-1-chat-app/backend/services/pb"
	"github.com/BrosSquad/ts-1-chat-app/backend/validators"
)

type UserMessage struct {
	User    models.User
	Message models.Message
}

type chatService struct {
	debugLogger      *logging.Debug
	errorLogger      *logging.Error
	validator        validators.Validator
	connections      sync.Map
	numOfConnections uint64
	messages         chan UserMessage
	db               *gorm.DB
	pb.UnimplementedChatServer
}

func New(db *gorm.DB, validator validators.Validator, debugLogger *logging.Debug, errorLogger *logging.Error, buffer uint16) pb.ChatServer {
	return &chatService{
		debugLogger:      debugLogger,
		errorLogger:      errorLogger,
		validator:        validator,
		connections:      sync.Map{},
		numOfConnections: 0,
		db:               db,
		messages:         make(chan UserMessage, buffer),
	}
}


