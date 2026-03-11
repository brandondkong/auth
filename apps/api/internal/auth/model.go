package auth

import (
	"time"

	"github.com/brandondkong/auth/internal/user"
	"github.com/brandondkong/auth/pkg/database"
	"github.com/google/uuid"
)

type RefreshToken struct {
	database.Model	`gorm:"embedded"`
	ExpiresAt	time.Time
	Revoked		*bool
	UserId		uuid.UUID
	User		user.User	`gorm:"foreignKey:UserId"`
	TokenId		string		`gorm:"uniqueIndex"`
}
