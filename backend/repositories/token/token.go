package token

import (
	"context"
	"errors"
	"github.com/BrosSquad/ts-1-chat-app/backend/models"
	"github.com/BrosSquad/ts-1-chat-app/backend/repositories"
	"gorm.io/gorm"

	"github.com/BrosSquad/ts-1-chat-app/backend/logging"
)

type (
	Repository interface {
		Find(context.Context, *[16]byte) (models.Token, error)
		Create(context.Context, models.Token) (models.Token, error)
		Delete(context.Context, string) error
	}

	repository struct {
		db          *gorm.DB
		errorLogger *logging.Error
	}
)

func (r repository) Create(ctx context.Context, token models.Token) (models.Token, error) {
	tx := r.db.WithContext(ctx)

	result := tx.Create(&token)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.Token{}, repositories.ErrNotFound
		}

		r.errorLogger.
			Err(result.Error).
			Str("type", "database").
			Str("query", result.Statement.SQL.String()).
			Msg("error while reading the database")

		return models.Token{}, repositories.DatabaseError{
			Inner:   result.Error,
			Message: "database error while fetching user",
		}
	}

	return token, nil
}

func (r repository) Delete(ctx context.Context, s string) error {
	panic("implement me")
}

func (r repository) Find(ctx context.Context, id *[16]byte) (models.Token, error) {
	panic("implement me")
}

func New(db *gorm.DB, errorLogger *logging.Error) Repository {
	return repository{
		db:          db,
		errorLogger: errorLogger,
	}
}
