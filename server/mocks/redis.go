package mocks

import (
	"errors"

	"github.com/michaelg9/ISOC/server/models"
)

// MockTokens implements the TokenControl interface so it can be used to mock the token back-end
type MockTokens struct{}

func (mTkns *MockTokens) CheckAccessToken(tokenString string) (models.User, error) {
	if tokenString == AccessToken {
		return models.User{ID: Users[0].ID, Admin: Users[0].Admin}, nil
	}
	return models.User{}, errors.New("Token is invalid.")
}

func (mTkns *MockTokens) CheckRefreshToken(tokenString string) (models.User, error) {
	if tokenString == RefreshToken {
		return models.User{ID: Users[0].ID, Admin: Users[0].Admin}, nil
	}
	return models.User{}, errors.New("Token is invalid.")
}

func (mTkns *MockTokens) NewAccessToken(user models.User) (string, error) {
	return AccessToken, nil
}

func (mTkns *MockTokens) NewRefreshToken(user models.User) (string, error) {
	return RefreshToken, nil
}

func (mTkns *MockTokens) InvalidateToken(tokenString string) error {
	if !(tokenString == RefreshToken || tokenString == AccessToken) {
		return errors.New("Token is invalid.")
	}
	return nil
}
