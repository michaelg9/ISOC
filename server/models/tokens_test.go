package models

import (
	"testing"
	"time"
)

func TestTokenGeneration(t *testing.T) {
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

	tokens := Tokens{}

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
