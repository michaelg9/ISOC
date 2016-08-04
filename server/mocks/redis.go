package mocks

import (
	"errors"
	"time"

	"github.com/michaelg9/ISOC/server/models"
)

// MockTokens implements the TokenControl interface so it can be used to mock the token back-end
type MockTokens struct{}

func (mTkns *MockTokens) CheckToken(tokenString string) (email string, err error) {
	if tokenString == JWT {
		return "user@usermail.com", nil
	}
	return "", errors.New("Token is invalid.")
}

func (mTkns *MockTokens) NewToken(user models.User, duration time.Duration) (string, error) {
	return JWT, nil
}

func (mTkns *MockTokens) InvalidateToken(tokenString string) error {
	if tokenString != JWT {
		return errors.New("Token already invalid.")
	}
	return nil
}
