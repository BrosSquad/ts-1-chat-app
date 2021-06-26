package models

type User struct {
	Model
	Name     string
	Surname  string
	Email    string `gorm:"uniqueIndex"`
	Password string

	Messages []Message `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignkey:UserID;"`
	Tokens   []Token `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignkey:UserID;"`
}
