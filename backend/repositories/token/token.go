package token

import (
	"context"
	"github.com/BrosSquad/ts-1-chat-app/backend/models"
	"gorm.io/gorm"

	"github.com/BrosSquad/ts-1-chat-app/backend/logging"
)

type (
	Repository interface {
		Find(context.Context, string) (models.Token, error)
		Create(context.Context, models.Token) (models.Token, error)
		Delete(context.Context, string) error
	}

	repository struct {
		db          *gorm.DB
		errorLogger *logging.Error
	}
)

func (r repository) Create(ctx context.Context, token models.Token) (models.Token, error) {
	panic("implement me")
}

func (r repository) Delete(ctx context.Context, s string) error {
	panic("implement me")
}

func (r repository) Find(ctx context.Context, id string) (models.Token, error) {
	panic("implement me")
}

func New(db *gorm.DB, errorLogger *logging.Error) Repository {
	return repository{
		db:          db,
		errorLogger: errorLogger,
	}
}
