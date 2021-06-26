package repositories

import "errors"

var (
	ErrAlreadyExists = errors.New("model already exists")
	ErrNotFound      = errors.New("not found")
	ErrDeleteFailed  = errors.New("delete failed")
)

type DatabaseError struct {
	Inner   error
	Message string
}

func (d DatabaseError) Error() string {
	return d.Message
}
