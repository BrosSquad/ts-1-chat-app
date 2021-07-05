package token

import (
	"context"
	"errors"
	"fmt"
	"github.com/BrosSquad/ts-1-chat-app/backend/models"
	"github.com/BrosSquad/ts-1-chat-app/backend/repositories"
	"gorm.io/gorm"

	"github.com/BrosSquad/ts-1-chat-app/backend/logging"
)

type (
	Repository interface {
		Find(context.Context, []byte) (models.Token, error)
		Create(context.Context, models.Token) (models.Token, error)
		Delete(context.Context, []byte) error
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

func (r repository) Delete(ctx context.Context, tokenId []byte) error {
	result := r.db.
		WithContext(ctx).
		Where("id = ?", tokenId).
		Limit(1).
		Delete(models.Token{})

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return repositories.ErrNotFound
		}

		r.errorLogger.
			Err(result.Error).
			Str("type", "database").
			Str("query", result.Statement.SQL.String()).
			Msg("error while reading the database")

		return repositories.DatabaseError{
			Inner:   result.Error,
			Message: "database error while deleting token",
		}
	}

	if result.RowsAffected != 1 {
		return repositories.DatabaseError{
			Inner:   result.Error,
			Message: fmt.Sprintf("Database error: %d rows affected", result.RowsAffected),
		}
	}

	return nil
}

func (r repository) Find(ctx context.Context, id []byte) (models.Token, error) {
	tx := r.db.WithContext(ctx)
	token := models.Token{
		ID: id,
	}

	result := tx.First(&token)

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

func New(db *gorm.DB, errorLogger *logging.Error) Repository {
	return repository{
		db:          db,
		errorLogger: errorLogger,
	}
}
