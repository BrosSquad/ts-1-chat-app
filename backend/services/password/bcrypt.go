package password

// TODO: Optimize for Buffer allocations

import (
	"errors"
	"github.com/BrosSquad/ts-1-chat-app/backend/logging"
	"golang.org/x/crypto/bcrypt"
)

type bcryptHasher struct {
	errorLogger *logging.Error
	cost        int
}

func (b bcryptHasher) Hash(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)

	if err != nil {
		b.errorLogger.
			Err(err).
			Int("cost", b.cost).
			Msg("Error while hashing password")

		return ""
	}

	return string(hash)
}

func (b bcryptHasher) Verify(hash string, plaintext string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plaintext))

	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrMismatchedHashAndPassword
		}

		return err
	}

	return nil
}

func NewBCryptHasher(errorLogger *logging.Error, cost int) Hasher {
	return bcryptHasher{
		cost:        cost,
		errorLogger: errorLogger,
	}
}
