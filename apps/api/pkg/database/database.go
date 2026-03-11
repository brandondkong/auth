package database

import (
	"errors"

	"github.com/brandondkong/auth/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func StartDatabase(configs config.Config) (*gorm.DB, error) {
	if db != nil {
		return nil, nil
	}
	database, err := gorm.Open(postgres.Open(configs.DatabaseUrl), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db = database

	return database, nil
}

func GetDatabase() (*gorm.DB, error) {
	if db == nil {
		return nil, errors.New("database not initialized")
	}
	return db, nil
}
