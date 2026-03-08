package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"time"

	"github.com/brandondkong/auth/internal/user"
	"github.com/brandondkong/auth/pkg/database"
	"gorm.io/gorm"
)

const numBytesInMagicLink uint = 32
const magicLinkExpiryDurationHours uint = 24

var ErrInvalidToken error = errors.New("invalid token")

func generateSecureToken() (string, error) {
	buffer := make([]byte, numBytesInMagicLink)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}
	token := base64.URLEncoding.EncodeToString(buffer)
	return token, nil 
}

func hashSecureToken(token string) ([]byte, error) {
	hash := sha256.New()
	hash.Write([]byte(token))

	return hash.Sum(nil), nil
}

func GenerateMagicLink(email string, request *http.Request) error {
	token, err := generateSecureToken()
	if err != nil {
		return err
	}

	hash, err := hashSecureToken(token)
	if err != nil {
		return err
	}

	databaseSafeToken := string(hash)
	
	magicLinkToken := MagicLinkToken{
		Email:		email,
		UserAgent: request.UserAgent(),
		IPAddress: request.Host,
		Token:		databaseSafeToken,
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(magicLinkExpiryDurationHours)),
	}

	db, err := database.GetDatabase()
	if err != nil {
		return err
	}

	err = db.Create(&magicLinkToken).Error
	if err != nil {
		return err
	}

	return nil
}

func ConsumeMagicLink(token string) (*user.User, error) {
	hash, err := hashSecureToken(token)
	if err != nil {
		return nil, err
	}

	databaseSafeToken := string(hash)

	// Query for hashed token against database
	db, err := database.GetDatabase()
	if err != nil {
		return nil, err
	}

	magicLinkToken := MagicLinkToken{}

	var existingUser *user.User

	// Wrap in a transaction
	err = db.Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&magicLinkToken).Where("token = ? AND used = false AND expires_at > ?", databaseSafeToken, time.Now()).Updates(map[string]any{
			"used": true,
			"used_at": time.Now(),
		})
		err := res.Error

		if err != nil {
			return err
		}

		if res.RowsAffected != 1 {
			return ErrInvalidToken
		}

		// Retrieve the email
		err = tx.Where("token = ?", databaseSafeToken).First(&magicLinkToken).Error
		if err != nil {
			return err
		}

		existingUser, err = user.GetUserByEmail(magicLinkToken.Email, tx)
		if err != nil {
			if errors.Is(err, user.ErrUserNotFound) {
				// Create the user
				existingUser, err = user.CreateUser(magicLinkToken.Email, tx)
				if err != nil {
					return err
				}
				return nil
	
			}
			return err
		}

		return nil
	})

	return existingUser, err
}
