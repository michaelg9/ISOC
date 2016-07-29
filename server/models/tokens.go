package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	hmacSecret = "secret"
)

// Tokens is the struct for outputting the two types of tokens: a long-lived refresh token
// and a short-lived access token.
type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// TokenControl is the interface for all functions related to tokens.
type TokenControl interface {
	CheckToken(token string) (email string, err error)
	NewToken(user User, duration time.Duration) (string, error)
	InvalidateToken(token string) error
}

// Tokenstore implements the TokenControl interface.
type Tokenstore struct{}

// NewTokenstore returns a new tokenstore struct.
func NewTokenstore() Tokenstore {
	return Tokenstore{}
}

// CheckToken checks if the given token is valid.
func (tkns *Tokenstore) CheckToken(tokenString string) (email string, err error) {
	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		err = errors.New("Token is invalid.")
		return
	}

	email, ok = claims["sub"].(string)
	if !ok || email == "" {
		err = errors.New("Claims do not contain user email.")
		return
	}

	return
}

// NewToken generates a new token for the given user and valid for the given duration.
func (tkns *Tokenstore) NewToken(user User, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Email,
		"exp": time.Now().Add(duration).Unix(), // Sets expiration time one day from now
	})

	return token.SignedString([]byte(hmacSecret))
}

// InvalidateToken stores to given token in Redis to blacklist it.
func (tkns *Tokenstore) InvalidateToken(token string) error {
	// TODO: Implement
	return errors.New("Not implemented.")
}

func keyFunc(token *jwt.Token) (interface{}, error) {
	// Check if token was signed with the right signing method
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}
	// If yes return the key
	return []byte(hmacSecret), nil
}
