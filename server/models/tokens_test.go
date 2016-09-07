package models

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testDBNR = 1

func setupRedis() *Tokenstore {
	return NewTokenstore(os.Getenv("REDIS_HOST")+":6379", testDBNR)
}

func cleanUpRedis(tokens *Tokenstore) {
	c := tokens.Get()
	defer c.Close()

	_, err := c.Do("FLUSHDB")
	if err != nil {
		panic(err)
	}

	tokens.Close()
}

func TestNewToken(t *testing.T) {
	tokens := setupRedis()
	defer cleanUpRedis(tokens)

	untilTomorrow := time.Hour * time.Duration(24)
	untilYesterday := time.Hour * time.Duration(-24)
	testRole := "testRole"
	var tests = []struct {
		id          int
		admin       bool
		duration    time.Duration
		role        string
		expectedErr error
	}{
		{1, true, untilTomorrow, testRole, nil},
		{1, false, untilTomorrow, testRole, nil},
		{1, true, untilYesterday, testRole, errors.New("Token is expired")},
		{0, true, untilTomorrow, testRole, errNoSubClaim},
		{1, true, untilTomorrow, "wrongRole", errWrongRole},
	}

	for _, test := range tests {
		user := User{ID: test.id, Admin: test.admin}
		token, err := newToken(user, test.duration, test.role)
		assert.Empty(t, err)

		result, err := tokens.checkToken(token, testRole)
		if test.expectedErr != nil {
			assert.EqualError(t, test.expectedErr, err.Error())
		} else if assert.NoError(t, err) {
			assert.Equal(t, user, result)
		}
	}
}

func TestInvalidateToken(t *testing.T) {
	tokens := setupRedis()
	defer cleanUpRedis(tokens)

	untilTomorrow := time.Hour * time.Duration(24)
	untilYesterday := time.Hour * time.Duration(-24)
	testUser := User{ID: 1, Admin: true}
	testRole := "testRole"

	var tests = []struct {
		duration                   time.Duration
		expectedErrForInvalidation error
		expectedErrForChecking     error
	}{
		{untilTomorrow, nil, errBlacklisted},
		{untilYesterday, errInvalid, errors.New("Token is expired")},
	}

	for _, test := range tests {
		token, err := newToken(testUser, test.duration, testRole)
		assert.Empty(t, err)

		err = tokens.InvalidateToken(token)
		if test.expectedErrForInvalidation != nil {
			assert.EqualError(t, test.expectedErrForInvalidation, err.Error())
		} else {
			assert.NoError(t, err)
		}

		_, err = tokens.checkToken(token, testRole)
		if test.expectedErrForChecking != nil {
			assert.EqualError(t, test.expectedErrForChecking, err.Error())
		} else {
			assert.NoError(t, err)
		}
	}
}
