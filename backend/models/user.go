package models

type User struct {
	Model
	Name     string
	Surname  string
	Email    string `gorm:"uniqueIndex"`
	Password string

	Messages []Message
}
