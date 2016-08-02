package models

import (
	"os"
	"testing"
	"time"
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

func TestTokenGeneration(t *testing.T) {
	tokens := setupRedis()
	defer cleanUpRedis(tokens)

	untilTomorrow := time.Hour * time.Duration(24)
	untilYesterday := time.Hour * time.Duration(-24)
	var tests = []struct {
		email                  string
		duration               time.Duration
		expectedErrForChecking string
	}{
		{"user@usermail.com", untilTomorrow, ""},
		{"user@usermail.com", untilYesterday, "Token is expired"},
		{"", untilTomorrow, "Claims do not contain user email."},
	}

	for _, test := range tests {
		token, err := tokens.NewToken(User{Email: test.email}, test.duration)
		if err != nil {
			t.Errorf("\n...got error = %v", err)
		}

		email, err := tokens.CheckToken(token)
		if err != nil {
			if err.Error() != test.expectedErrForChecking {
				t.Errorf("\n...expected = %v\n...obtained = %v", test.expectedErrForChecking, err)
			}
		}
		if email != test.email && test.expectedErrForChecking == "" {
			t.Errorf("\n...expected = %v\n...obtained = %v", test.email, email)
		}
	}
}

func TestInvalidateToken(t *testing.T) {
	tokens := setupRedis()
	defer cleanUpRedis(tokens)

	untilTomorrow := time.Hour * time.Duration(24)
	untilYesterday := time.Hour * time.Duration(-24)

	var tests = []struct {
		duration                   time.Duration
		expectedErrForInvalidation string
		expectedErrForChecking     string
	}{
		{untilTomorrow, "", "Token is blacklisted."},
		{untilYesterday, "Token already invalid.", "Token is expired"},
	}

	for _, test := range tests {
		token, err := tokens.NewToken(User{Email: "user@usermail.com"}, test.duration)
		if err != nil {
			t.Errorf("\n...got error = %v", err)
		}

		err = tokens.InvalidateToken(token)
		if err != nil {
			if err.Error() != test.expectedErrForInvalidation {
				t.Errorf("\n...expected = %v\n...obtained = %v", test.expectedErrForInvalidation, err)
			}
		}

		_, err = tokens.CheckToken(token)
		if err != nil {
			if err.Error() != test.expectedErrForChecking {
				t.Errorf("\n...expected = %v\n...obtained = %v", test.expectedErrForChecking, err)
			}
		}
	}
}
