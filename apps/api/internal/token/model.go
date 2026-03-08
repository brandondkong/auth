package token

import (
	"time"

	"github.com/brandondkong/auth/pkg/database"
)

type MagicLinkToken struct {
	database.Model
	Email		string
	Token		string		`gorm:"primaryKey"`
	IPAddress	string
	UserAgent	string
	Used		bool
	UsedAt		time.Time
	ExpiresAt	time.Time
}
