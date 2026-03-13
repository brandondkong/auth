package jwtutil

import "github.com/golang-jwt/jwt/v5"

func ParseToken(tkn string, key string) (*jwt.Token, error) {
	parsed, err := jwt.ParseWithClaims(tkn, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {

		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}
	
	return parsed, nil
}
