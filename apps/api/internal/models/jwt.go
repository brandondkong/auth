package models

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	TokenId		string		`gorm:"primaryKey"`
	ExpiresAt	time.Time
	Revoked		*bool
	UserId		uuid.UUID
}

