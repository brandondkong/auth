package user

import (
	"errors"

	"github.com/brandondkong/auth/internal/models"
	"github.com/brandondkong/auth/pkg/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetUserById(id uuid.UUID, tx *gorm.DB) (*models.User, error) {
	var db *gorm.DB = tx
	var err error

	if tx == nil {
		db, err = database.GetDatabase()
		if err != nil {
			return nil, err
		}
	}
	
	var existing models.User
	if err = db.Where("id = ?", id).First(&existing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &existing, nil
}

func GetUserByEmail(email string, tx *gorm.DB) (*models.User, error) {
	var db *gorm.DB = tx
	var err error

	if tx == nil {
		db, err = database.GetDatabase()
		if err != nil {
			return nil, err
		}
	}
	
	var existing models.User
	if err = db.Where("email = ?", email).First(&existing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &existing, nil
}

func CreateUser(email string, tx *gorm.DB) (*models.User, error) {
	var db *gorm.DB = tx
	var err error

	if tx == nil {
		db, err = database.GetDatabase()
		if err != nil {
			return nil, err
		}
	}
	
	// TODO: validate email

	user := &models.User{
		Email:		email,
	}

	err = db.Create(user).Error
	return user, err
}
