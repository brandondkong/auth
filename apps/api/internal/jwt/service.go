package jwt

import (
	"log"
	"time"

	"github.com/brandondkong/auth/internal/config"
	"github.com/brandondkong/auth/internal/user"
	"github.com/brandondkong/auth/pkg/cryptoutil"
	"github.com/brandondkong/auth/pkg/database"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateTokenPair(user *user.User) (*TokenPair, error) {
	config, err := config.LoadConfigs()
	if err != nil {
		return nil, err
	}

	refreshSignature, _, err := CreateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	accessSignature, _, err := CreateJwtToken(user, config.JwtAccessSigningKey, time.Now().Add(AccessTokenLifeTime))
	if err != nil {
		return nil, err
	}

	var tokenPair TokenPair = TokenPair{
		Refresh: refreshSignature,
		Access: accessSignature,
	}

	return &tokenPair, nil
}

func CreateJwtToken(user *user.User, key string, expires time.Time) (string, uuid.UUID, error) {
	tokenId := uuid.New()
	expiresAt := jwt.NewNumericDate(expires)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
			jwt.RegisteredClaims{
				ExpiresAt: expiresAt,			
				IssuedAt: jwt.NewNumericDate(time.Now()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				Issuer: TokenIssuer,
				Subject: user.ID.String(),
				ID: tokenId.String(),
			})

	signed, err := token.SignedString([]byte(key))
	if err != nil {
		return "", uuid.Nil, err
	}

	return signed, tokenId, nil
}

func CreateRefreshToken(user *user.User) (string, uuid.UUID, error) {
	config, err := config.LoadConfigs()
	if err != nil {
		return "", uuid.Nil, err
	}

	expiration := time.Now().Add(RefreshTokenLifeTime)
	signature, tokenId, err := CreateJwtToken(user, config.JwtRefreshSigningKey, expiration) 
	if err != nil {
		return "", uuid.Nil, err
	}

	// Save this into the database
	db, err := database.GetDatabase()
	if err != nil {
		return "", uuid.Nil, err
	}
	
	// Hash the token ID
	hashedId, err := cryptoutil.HashString(tokenId.String())
	if err != nil {
		return "", uuid.Nil, err
	}

	var refreshToken RefreshToken = RefreshToken{
		ExpiresAt: expiration,
		UserId: user.ID,
		TokenId: string(hashedId),
	}

	err = db.Create(&refreshToken).Error
	if err != nil {
		return "", uuid.Nil, err
	}
	
	return signature, tokenId, nil
}

func RotateTokens(tkn string) (*TokenPair, error) {
	parsed, err := jwt.ParseWithClaims(tkn, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		configs, err := config.LoadConfigs()
		if err != nil {
			return nil, err
		}

		return []byte(configs.JwtRefreshSigningKey), nil
	})
	
	if err != nil {
		return nil, err
	} else if claims, ok := parsed.Claims.(*jwt.RegisteredClaims); ok {
		// OKAY
		// From the JTI,
		hashed, err := cryptoutil.HashString(claims.ID)
		if err != nil {
			return nil, err
		}

		db, err := database.GetDatabase()
		if err != nil {
			return nil, err
		}

		var refresh RefreshToken

		err = db.Transaction(func(tx *gorm.DB) error {
			res := tx.Model(&refresh).Where("token_id = ? AND expires_at > ? AND revoked = false", hashed, time.Now()).Updates(map[string]any{
				"revoked":	true,
			})
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected != 1 {
				return ErrInvalidRefreshToken
			}

			// Find
			err := tx.Preload("User").Where("token_id = ?", hashed).Find(&refresh).Error
			if err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return nil, err
		}

		// At this point, refresh is valid
		// Create a new TokenPair and return it
		tokenPair, err := CreateTokenPair(&refresh.User)
		if err != nil {
			return nil, err
		}
		return tokenPair, nil
		
	} else {
		return nil, ErrUnknownClaimsType
	}
}

func CleanupStaleRefreshTokens() error {
	db, err := database.GetDatabase()
	if err != nil {
		return err
	}

	res := db.Unscoped().Where("expires_at < ? OR revoked = true", time.Now()).Delete(&RefreshToken{})
	if res.Error != nil {
		return res.Error
	}
	
	log.Printf("cleaned up %d stale refresh tokens\n", res.RowsAffected)
	return nil
}

