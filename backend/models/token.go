package models

import "time"

type Token struct {
	ID   []byte `gorm:"primary"`
	Hash []byte

	CreatedAt time.Time

	UserID uint64 `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User   User
}
