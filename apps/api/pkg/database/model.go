package database

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Migratable interface {
	TableName() string
}

type Model struct {
	ID			uuid.UUID		`gorm:"primaryKey"`
	CreatedAt 	time.Time		`gorm:"autoCreateTime"`
	UpdatedAt 	time.Time		`gorm:"autoUpdateTime:milli"`
	DeletedAt 	gorm.DeletedAt 	`gorm:"index"`
}

func (m *Model) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	m.CreatedAt = time.Now()
	return
}
