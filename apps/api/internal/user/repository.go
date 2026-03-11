package user

import (
	"errors"

	"github.com/brandondkong/auth/pkg/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrUserNotFound error = errors.New("user not found")

func GetUserById(id uuid.UUID, tx *gorm.DB) (*User, error) {
	var db *gorm.DB = tx
	var err error

	if tx == nil {
		db, err = database.GetDatabase()
		if err != nil {
			return nil, err
		}
	}
	
	var existing User
	if err = db.Where("id = ?", id).First(&existing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &existing, nil
}

func GetUserByEmail(email string, tx *gorm.DB) (*User, error) {
	var db *gorm.DB = tx
	var err error

	if tx == nil {
		db, err = database.GetDatabase()
		if err != nil {
			return nil, err
		}
	}
	
	var existing User
	if err = db.Where("email = ?", email).First(&existing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &existing, nil
}

func CreateUser(email string, tx *gorm.DB) (*User, error) {
	var db *gorm.DB = tx
	var err error

	if tx == nil {
		db, err = database.GetDatabase()
		if err != nil {
			return nil, err
		}
	}
	
	// TODO: validate email

	user := &User{
		Email:		email,
	}

	err = db.Create(user).Error
	return user, err
}
