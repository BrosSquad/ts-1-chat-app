package chat

import (
	"sync/atomic"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/BrosSquad/ts-1-chat-app/backend/models"
	"github.com/BrosSquad/ts-1-chat-app/backend/services/pb"
)

func (c *chatService) Connect(req *pb.ConnectRequest, client pb.Chat_ConnectServer) error {
	value := atomic.AddUint64(&c.numOfConnections, 1)
	c.debugLogger.Debug().
		Str("type", "connections").
		Uint64("numberOfConnections", value).
		Msg("Number of concurrent connections")

	c.connections.Store(req.UserId, client)
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
		_ = client.Send(&pb.MessageResponse{
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
				c.connections.Range(func(key, value interface{}) bool {
					client := value.(pb.Chat_ConnectServer)
					_ = client.Send(&pb.MessageResponse{
						User: &pb.User{
							Id:      um.User.ID,
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
			value := atomic.AddUint64(&c.numOfConnections, ^uint64(0))
			c.connections.Delete(req.UserId)

			c.debugLogger.Debug().
				Str("type", "connections").
				Uint64("numberOfConnections", value).
				Msg("Number of concurrent connections")

			return nil
		}
	}
}
