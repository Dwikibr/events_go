package utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

func GenerateToken(username string, id int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"id":       id,
		"exp":      time.Now().Add(time.Hour * 2).Unix(),
	})

	return token.SignedString([]byte(GetSignedKey()))
}
