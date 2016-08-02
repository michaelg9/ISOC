package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/garyburd/redigo/redis"
)

const (
	hmacSecret = "secret"
)

// Tokens is the struct for outputting the two types of tokens: a long-lived refresh token
// and a short-lived access token.
type Tokens struct {
	AccessToken  string `json:"accessToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

// TokenControl is the interface for all functions related to tokens.
type TokenControl interface {
	CheckToken(token string) (email string, err error)
	NewToken(user User, duration time.Duration) (string, error)
	InvalidateToken(tokenString string) error
}

// Tokenstore implements the TokenControl interface.
type Tokenstore struct {
	*redis.Pool
}

// NewTokenstore returns a new tokenstore struct.
func NewTokenstore(server string, dbNR int) *Tokenstore {
	pool := redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", server)
		if err != nil {
			return nil, err
		}
		if _, err := c.Do("SELECT", dbNR); err != nil {
			c.Close()
			return nil, err
		}
		// TODO: Specify redis password
		return c, nil
	}, 3) // 3 is the max number of idle connections
	return &Tokenstore{pool}
}

// CheckToken checks if the given token is valid.
func (tkns *Tokenstore) CheckToken(tokenString string) (email string, err error) {
	if isInvalidated(tkns.Pool, tokenString) {
		err = errors.New("Token is blacklisted.")
		return
	}

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
func (tkns *Tokenstore) InvalidateToken(tokenString string) error {
	// Get redis connection
	c := tkns.Pool.Get()
	defer c.Close()

	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		return errors.New("Token already invalid.")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		return errors.New("Token already invalid.")
	}

	remainingValidity := getTokenRemainingValidity(claims["exp"])

	_, err = c.Do("SET", tokenString, tokenString)
	if err == nil {
		_, err = c.Do("EXPIRE", tokenString, remainingValidity)
	}

	return err
}

func isInvalidated(pool *redis.Pool, tokenString string) bool {
	c := pool.Get()
	redisToken, _ := c.Do("GET", tokenString)
	return redisToken != nil
}

// Function from brainattica
func getTokenRemainingValidity(timestamp interface{}) int {
	// TODO: Add offset
	if validity, ok := timestamp.(float64); ok {
		tm := time.Unix(int64(validity), 0)
		remainer := tm.Sub(time.Now())
		if remainer > 0 {
			return int(remainer.Seconds())
		}
	}
	return 0
}

func keyFunc(token *jwt.Token) (interface{}, error) {
	// Check if token was signed with the right signing method
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}
	// If yes return the key
	return []byte(hmacSecret), nil
}
