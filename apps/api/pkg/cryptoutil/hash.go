package cryptoutil

import (
	"crypto/sha256"
	"encoding/base64"
)

func HashString(str string) (string, error) {
	hash := sha256.New()
	hash.Write([]byte(str))

	return base64.RawURLEncoding.EncodeToString(hash.Sum(nil)), nil
}

