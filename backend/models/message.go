package models

type Message struct {
	Model
	Text   string
	UserID uint64

	User User
}
