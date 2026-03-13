package models

import (
	"github.com/brandondkong/auth/pkg/database"
	"github.com/google/uuid"
)

type User struct {
	database.Model	`gorm:"embedded"`
	Email	string	`gorm:"unique;not null" json:"email"`
	Name	*string `json:"name"`
}

func (u User) TableName() string {
	return "users"
}

type OAuthAccount struct {
	User		User	`gorm:"foreignKey:UserId"`
	UserId		uuid.UUID	`gorm:"primaryKey"`
	Provider	string		`gorm:"primaryKey"`
	ProviderId	string
}

func (a OAuthAccount) TableName() string {
	return "oauth_accounts"
}

