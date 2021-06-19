package models

type User struct {
	Model
	Username string

	Messages []Message
}
