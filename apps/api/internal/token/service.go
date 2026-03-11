package token

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/brandondkong/auth/internal/user"
	"github.com/brandondkong/auth/pkg/cryptoutil"
	"github.com/brandondkong/auth/pkg/database"
	"gorm.io/gorm"
)

func generateSecureToken() (string, error) {
	buffer := make([]byte, NumBytesInMagicLink)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}
	token := base64.RawURLEncoding.EncodeToString(buffer)
	return token, nil 
}

func GenerateMagicLink(email string, request *http.Request) (string, error) {
	token, err := generateSecureToken()
	if err != nil {
		return "", err
	}

	databaseSafeToken, err := cryptoutil.HashString(token)
	if err != nil {
		return "", err
	}

	magicLinkToken := MagicLinkToken{
		Email:		email,
		UserAgent: request.UserAgent(),
		IPAddress: request.Host,
		Token:		databaseSafeToken,
		ExpiresAt: time.Now().Add(MagicLinkExpiryDuration),
	}

	db, err := database.GetDatabase()
	if err != nil {
		return "", err
	}

	err = db.Create(&magicLinkToken).Error
	if err != nil {
		return "", err
	}

	return token, nil
}

func ConsumeMagicLink(token string) (*user.User, error) {
	databaseSafeToken, err := cryptoutil.HashString(token)
	if err != nil {
		return nil, err
	}

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

func CleanupStaleTokens() error {
	db, err := database.GetDatabase()
	if err != nil {
		return err
	}

	res := db.Unscoped().Where("expires_at < ? OR used = true", time.Now()).Delete(&MagicLinkToken{})
	if res.Error != nil {
		return res.Error
	}
	
	log.Printf("cleaned up %d stale MagicLinkTokens\n", res.RowsAffected)
	return nil
}
