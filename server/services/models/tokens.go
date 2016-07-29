package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/michaelg9/ISOC/server/services/models"
)

// TokenControl is the interface for all functions related to tokens.
type TokenControl interface {
	CheckToken(token string) (email string, err error)
	NewToken(user User, duration time.Duration) (string, error)
	InvalidateToken(token string) error
}

// CheckToken checks if the given token is valid.
func CheckToken(token string) (email string, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(hmacSecret), nil
	})
	if err != nil {
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		// TODO: More meaningful error
		err = errors.New("Token not valid.")
		return
	}

	email = claims["sub"].(string)
	// TODO: Check if email exists
	return
}

// NewToken generates a new token for the given user and valid for the given duration.
func NewToken(user models.User, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Email,
		"exp": time.Now().Add(duration).Unix(), // Sets expiration time one day from now
	})

	return token.SignedString([]byte(hmacSecret))
}
