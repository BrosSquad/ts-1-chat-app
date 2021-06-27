package auth

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"hash"
	"strings"

	"google.golang.org/grpc/metadata"

	"github.com/BrosSquad/ts-1-chat-app/backend/models"
	"github.com/BrosSquad/ts-1-chat-app/backend/repositories/token"
	"github.com/BrosSquad/ts-1-chat-app/backend/utils"
)

const (
	IdSize    = 16
	ValueSize = 32
)

type (
	TokenService interface {
		Generate(context.Context, models.User) (string, error)
		Verify(context.Context, string, string) error
		Delete(context.Context, string) error
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

func (t tokenService) Delete(ctx context.Context, token string) error {
	id, _, err := t.parse(token)

	if err != nil {
		return err
	}

	return t.repo.Delete(ctx, id[:])
}

func (t tokenService) parse(token string) ([IdSize]byte, [ValueSize]byte, error) {
	keyAndValue := strings.SplitN(token, ".", 2)

	if len(keyAndValue) != 2 {
		return [16]byte{}, [32]byte{}, errors.New("token is invalid")
	}

	var id [IdSize]byte
	var value [ValueSize]byte

	_, err := base64.RawURLEncoding.Decode(id[:], utils.UnsafeBytes(keyAndValue[0]))

	if err != nil {
		return [16]byte{}, [32]byte{}, err
	}

	_, err = base64.RawURLEncoding.Decode(value[:], utils.UnsafeBytes(keyAndValue[1]))

	if err != nil {
		return [16]byte{}, [32]byte{}, err
	}

	return id, value, nil
}

func (t tokenService) Verify(ctx context.Context, tokenType, token string) error {
	id, value, err := t.parse(token)

	if err != nil {
		return err
	}

	tokenModel, err := t.repo.Find(ctx, id[:])

	if err != nil {
		return err
	}

	if tokenModel.Type != tokenType {
		return errors.New("invalid token type")
	}

	hashedValue := t.hash.Sum(value[:])

	if subtle.ConstantTimeCompare(hashedValue, tokenModel.Hash) == 0 {
		return errors.New("values are not equal")
	}

	return nil
}

func ExtractToken(ctx context.Context) (string, string, error) {
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return "", "", errors.New("metadata not found")
	}

	header := md.Get("authorization")

	if len(header) != 1 {
		return "", "", errors.New("authorization header not found")
	}

	typeAndToken := strings.SplitN(header[0], " ", 2)

	if len(typeAndToken) != 2 {
		return "", "", errors.New("authorization token invalid format: $TOKEN_TYPE$ $TOKEN$")
	}

	return typeAndToken[0], typeAndToken[1], nil
}

func NewTokenService(repo token.Repository, hash hash.Hash) TokenService {
	return &tokenService{
		repo: repo,
		hash: hash,
	}
}
