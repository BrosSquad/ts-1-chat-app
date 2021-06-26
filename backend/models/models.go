package models

import "time"

type Model struct {
	ID        uint64
	CreatedAt time.Time
}

func GetModels() []interface{} {
	return []interface{}{
		&User{},
		&Token{},
		&Message{},
	}
}
