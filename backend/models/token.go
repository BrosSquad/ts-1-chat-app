package models

import "time"

type Token struct {
	ID   []byte `gorm:"primary"`
	Hash []byte
	Type string `gorm:"index;default:bearer"`

	CreatedAt time.Time

	UserID uint64 `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User   User
}
