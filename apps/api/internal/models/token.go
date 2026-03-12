package models

import (
	"time"
)

type MagicLinkToken struct {
	Email		string
	Token		string		`gorm:"primaryKey"`
	IPAddress	string
	UserAgent	string
	Used		bool
	UsedAt		time.Time
	ExpiresAt	time.Time
}

