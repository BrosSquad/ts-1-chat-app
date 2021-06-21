package models

type User struct {
	Model
	Username string `gorm:"uniqueIndex"`

	Messages []Message
}
