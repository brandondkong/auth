package user

import (
	"github.com/brandondkong/auth/pkg/database"
)

type User struct {
	database.Model
	Email	string
	Name	*string
}

func (u User) TableName() string {
	return "user"
}
