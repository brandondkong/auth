package user

import (
	"github.com/brandondkong/auth/pkg/database"
)

type User struct {
	database.Model	`gorm:"embedded"`
	Email	string	`gorm:"unique;not null"`
	Name	*string
}

func (u User) TableName() string {
	return "user"
}
