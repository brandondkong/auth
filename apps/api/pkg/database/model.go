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
	ID			uuid.UUID		`gorm:"primaryKey" json:"id"`
	CreatedAt 	time.Time		`gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt 	time.Time		`gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt 	gorm.DeletedAt 	`gorm:"index" json:"deleted_at"`
}

func (m *Model) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	m.CreatedAt = time.Now()
	return
}
