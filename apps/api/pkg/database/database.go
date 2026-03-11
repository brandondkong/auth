package database

import (
	"errors"

	"github.com/brandondkong/auth/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func StartDatabase() (*gorm.DB, error) {
	if db != nil {
		return db, nil
	}
	configs, err := config.LoadConfigs()
	if err != nil {
		return nil, err
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

func Migrate(models ...any) error {
	db, err := GetDatabase()
	if err != nil {
		return err
	}

	err = db.AutoMigrate(models...)
	if err != nil {
		return err
	}

	return nil
}
