package chat

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/BrosSquad/ts-1-chat-app/backend/logging"
	"github.com/BrosSquad/ts-1-chat-app/backend/models"
	"github.com/BrosSquad/ts-1-chat-app/backend/services/pb"
)

type UserMessage struct {
	User    models.User
	Message models.Message
}

type chatService struct {
	debugLogger      *logging.Debug
	errorLogger      *logging.Error
	conns            sync.Map
	numOfconnections uint64
	messages         chan UserMessage
	db               *gorm.DB
	pb.UnimplementedChatServer
}

func New(db *gorm.DB, debugLogger *logging.Debug, errorLogger *logging.Error, buffer uint16) pb.ChatServer {
	return &chatService{
		debugLogger:      debugLogger,
		errorLogger:      errorLogger,
		conns:            sync.Map{},
		numOfconnections: 0,
		db:               db,
		messages:         make(chan UserMessage, buffer),
	}
}

func (c *chatService) SendMessage(ctx context.Context, in *pb.MessageRequest) (*pb.Empty, error) {
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

	return &pb.Empty{}, nil
}

func (c *chatService) Connect(req *pb.ConnectRequest, client pb.Chat_ConnectServer) error {
	value := atomic.AddUint64(&c.numOfconnections, 1)
	c.debugLogger.Debug().
		Str("type", "connections").
		Uint64("numberOfConnections", value).
		Msg("Number of conccurrent connections")

	c.conns.Store(req.UserId, client)
	messages := make([]models.Message, 0, 50)

	result := c.db.WithContext(client.Context()).
		Model(&models.Message{}).
		Preload("User").
		Order("created_at DESC").
		Limit(50).
		Find(&messages)

	if result.Error != nil {
		c.errorLogger.
			Err(result.Error).
			Str("type", "database").
			Str("query", result.Statement.SQL.String()).
			Msg("Error while fetching latest messages")

		return status.Error(codes.Internal, "Error while fetching latest messages")
	}

	for _, message := range messages {
		client.Send(&pb.MessageResponse{
			User: &pb.User{
				Id:      message.User.ID,
				Email:   message.User.Email,
				Name:    message.User.Name,
				Surname: message.User.Surname,
			},
			Text:      message.Text,
			CreatedAt: message.CreatedAt.Format(time.RFC3339),
		})
	}

	for {
		select {
		case um := <-c.messages:
			go func(um UserMessage) {
				c.conns.Range(func(key, value interface{}) bool {
					client := value.(pb.Chat_ConnectServer)
					client.Send(&pb.MessageResponse{
						User: &pb.User{
							Id:       um.User.ID,
							Email:   um.User.Email,
							Name:    um.User.Name,
							Surname: um.User.Surname,
						},
						Text:      um.Message.Text,
						CreatedAt: um.Message.CreatedAt.Format(time.RFC3339),
					})

					return true
				})
			}(um)
		case <-client.Context().Done():
			value := atomic.AddUint64(&c.numOfconnections, ^uint64(0))
			c.conns.Delete(req.UserId)

			c.debugLogger.Debug().
				Str("type", "connections").
				Uint64("numberOfConnections", value).
				Msg("Number of conccurrent connections")

			return nil
		}
	}
}
