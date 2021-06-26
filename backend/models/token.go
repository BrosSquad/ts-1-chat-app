package models

import "time"

type Token struct {
	ID   string `gorm:"primary"`
	Hash string

	CreatedAt time.Time

	UserID uint64 `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User User
}
