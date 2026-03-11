package jwt

import (
	"time"

	"github.com/brandondkong/auth/internal/user"
	"github.com/google/uuid"
)

type RefreshToken struct {
	TokenId		string		`gorm:"primaryKey"`
	ExpiresAt	time.Time
	Revoked		*bool
	UserId		uuid.UUID
	User		user.User	`gorm:"foreignKey:UserId"`
}

type TokenPair struct {
	Refresh		string	`json:"-"`
	Access		string	`json:"access"`
}
