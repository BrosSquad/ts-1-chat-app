package password

import "errors"

var (
	ErrMismatchedHashAndPassword = errors.New("password does not match the hash")
)

type Hasher interface {
	Hash(string) string
	Verify(hash string, plaintext string) error
}
