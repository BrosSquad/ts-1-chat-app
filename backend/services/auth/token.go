package auth

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/BrosSquad/ts-1-chat-app/backend/utils"
	"hash"
	"strings"

	"github.com/BrosSquad/ts-1-chat-app/backend/models"
	"github.com/BrosSquad/ts-1-chat-app/backend/repositories/token"
)

const (
	IdSize    = 16
	ValueSize = 32
)

type (
	TokenService interface {
		Generate(context.Context, models.User) (string, error)
		Verify(context.Context, string) error
	}

	tokenService struct {
		repo token.Repository
		hash hash.Hash
	}
)

func (t tokenService) Generate(ctx context.Context, user models.User) (string, error) {
	id := make([]byte, IdSize)
	value := make([]byte, ValueSize)

	_, err := rand.Read(id)

	if err != nil {
		return "", err
	}

	_, err = rand.Read(value)

	if err != nil {
		return "", err
	}

	_, err = t.repo.Create(ctx, models.Token{
		ID:     id,
		Hash:   t.hash.Sum(value[:]),
		UserID: user.ID,
	})

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s.%s",
		base64.RawURLEncoding.EncodeToString(id),
		base64.RawURLEncoding.EncodeToString(value),
	), nil
}

func (t tokenService) Verify(ctx context.Context, token string) error {
	keyAndValue := strings.SplitN(token, ".", 1)

	if len(keyAndValue) != 2 {
		return errors.New("token is invalid")
	}

	var id [IdSize]byte
	var value [ValueSize]byte

	_, err := base64.RawURLEncoding.Decode(id[:], utils.UnsafeBytes(keyAndValue[0]))

	if err != nil {
		return err
	}

	_, err = base64.RawURLEncoding.Decode(value[:], utils.UnsafeBytes(keyAndValue[1]))

	if err != nil {
		return err
	}

	tokenModel, err := t.repo.Find(ctx, &id)

	if err != nil {
		return err
	}

	hashedValue := t.hash.Sum(value[:])

	if subtle.ConstantTimeCompare(hashedValue, tokenModel.Hash) == 0 {
		return errors.New("values are not equal")
	}

	return nil
}

func NewTokenService(repo token.Repository, hash hash.Hash) TokenService {
	return &tokenService{
		repo: repo,
		hash: hash,
	}
}
