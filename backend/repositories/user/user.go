package user

import (
	"context"
	"errors"
	"gorm.io/gorm"

	"github.com/BrosSquad/ts-1-chat-app/backend/logging"
	"github.com/BrosSquad/ts-1-chat-app/backend/models"
	"github.com/BrosSquad/ts-1-chat-app/backend/repositories"
)

type (
	Repository interface {
		Find(ctx context.Context, id uint64) (models.User, error)
		FindByEmail(ctx context.Context, email string) (models.User, error)
		ExistsByEmail(ctx context.Context, email string) error
		Create(ctx context.Context, req models.User) (models.User, error)
		Delete(ctx context.Context, id uint64) error
	}

	repository struct {
		db          *gorm.DB
		errorLogger *logging.Error
	}
)

func (r repository) Find(ctx context.Context, id uint64) (models.User, error) {
	tx := r.db.WithContext(ctx)

	var user models.User

	result := tx.Where("id = ?", id).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.User{}, repositories.ErrNotFound
		}

		r.errorLogger.
			Err(result.Error).
			Str("type", "database").
			Str("query", result.Statement.SQL.String()).
			Msg("error while reading the database")

		return models.User{}, repositories.DatabaseError{
			Inner:   result.Error,
			Message: "database error while fetching user",
		}
	}

	return user, nil
}

func (r repository) FindByEmail(ctx context.Context, email string) (models.User, error) {
	tx := r.db.WithContext(ctx)

	var user models.User

	result := tx.Where("email = ?", email).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.User{}, repositories.ErrNotFound
		}

		r.errorLogger.
			Err(result.Error).
			Str("type", "database").
			Str("query", result.Statement.SQL.String()).
			Msg("error while reading the database")

		return models.User{}, repositories.DatabaseError{
			Inner:   result.Error,
			Message: "database error while fetching user",
		}
	}

	return user, nil
}

func (r repository) ExistsByEmail(ctx context.Context, email string) error {
	tx := r.db.WithContext(ctx)

	var user models.User

	result := tx.Where("email = ?", email).First(&user)

	// Invalid Database query
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		r.errorLogger.
			Err(result.Error).
			Str("type", "database").
			Str("query", result.Statement.SQL.String()).
			Msg("error while reading the database")

		return repositories.DatabaseError{
			Inner:   result.Error,
			Message: "database error while fetching user",
		}
	}

	// Record already exists
	if result.RowsAffected > 0 {
		return repositories.ErrAlreadyExists
	}

	return nil
}

func (r repository) Create(ctx context.Context, req models.User) (models.User, error) {
	tx := r.db.WithContext(ctx)
	result := tx.Create(&req)

	if result.Error != nil {
		r.errorLogger.
			Err(result.Error).
			Str("type", "database").
			Str("query", result.Statement.SQL.String()).
			Msg("cannot create new user")

		return models.User{}, repositories.DatabaseError{
			Inner:   result.Error,
			Message: "database error while inserting user",
		}
	}

	return req, nil
}

func (r repository) Delete(ctx context.Context, id uint64) error {
	tx := r.db.WithContext(ctx)

	result := tx.Delete(&models.User{}, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return repositories.ErrNotFound
		}

		return repositories.DatabaseError{
			Inner:   result.Error,
			Message: "delete failed for user model",
		}
	}

	if result.RowsAffected != 1 {
		return repositories.ErrDeleteFailed
	}

	return nil
}

func New(db *gorm.DB, logger *logging.Error) Repository {
	return repository{
		db:          db,
		errorLogger: logger,
	}
}
