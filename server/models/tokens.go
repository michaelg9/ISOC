package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/garyburd/redigo/redis"
)

var (
	errWrongRole    = errors.New("Wrong token role.")
	errNoRoleClaim  = errors.New("Claims do not contain role.")
	errNoAdminClaim = errors.New("Claims do not contain admin.")
	errNoSubClaim   = errors.New("Claims do not contain sub.")
	errBlacklisted  = errors.New("Token is blacklisted.")
	errInvalid      = errors.New("Token is invalid.")
)

const (
	maxIdleConns = 3
	hmacSecret   = "secret"

	accessTokenRole   = "access"
	accessTokenDelta  = time.Minute * time.Duration(10) // Ten minutes
	refreshTokenRole  = "refresh"
	refreshTokenDelta = time.Hour * time.Duration(168) // One week
)

// Tokens is the struct for outputting the two types of tokens: a long-lived refresh token
// and a short-lived access token.
type Tokens struct {
	AccessToken  string `json:"accessToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

// TokenControl is the interface for all functions related to tokens.
type TokenControl interface {
	CheckAccessToken(token string) (user User, err error)
	CheckRefreshToken(token string) (user User, err error)
	NewAccessToken(user User) (string, error)
	NewRefreshToken(user User) (string, error)
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
	}, maxIdleConns)
	return &Tokenstore{pool}
}

// CheckAccessToken validates the given access token. It returns an error if the
// given token is expired, does not contain the claims "sub", "admin" and "role", or
// when the role is something other than an access token.
func (tkns *Tokenstore) CheckAccessToken(token string) (User, error) {
	return tkns.checkToken(token, accessTokenRole)
}

// CheckRefreshToken validates the given access token. It returns an error if the
// given token is expired, blacklisted , does not contain the claims "sub", "admin"
// and "role", or  when the role is something other than an refresh token.
func (tkns *Tokenstore) CheckRefreshToken(token string) (User, error) {
	return tkns.checkToken(token, refreshTokenRole)
}

// NewAccessToken generates a new access token for the given user.
func (tkns *Tokenstore) NewAccessToken(user User) (string, error) {
	return newToken(user, accessTokenDelta, accessTokenRole)
}

// NewRefreshToken generates a new access token for the given user.
func (tkns *Tokenstore) NewRefreshToken(user User) (string, error) {
	return newToken(user, refreshTokenDelta, refreshTokenRole)
}

// InvalidateToken stores to given token in Redis to blacklist it.
func (tkns *Tokenstore) InvalidateToken(tokenString string) error {
	// Get redis connection
	c := tkns.Pool.Get()
	defer c.Close()

	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		return errInvalid
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		return errInvalid
	}

	remainingValidity := getTokenRemainingValidity(claims["exp"])

	_, err = c.Do("SET", tokenString, tokenString)
	if err == nil {
		_, err = c.Do("EXPIRE", tokenString, remainingValidity)
	}

	return err
}

// checkToken checks if the given token is valid.
func (tkns *Tokenstore) checkToken(tokenString, tokenRole string) (user User, err error) {
	if isInvalidated(tkns.Pool, tokenString) {
		err = errBlacklisted
		return
	}

	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		err = errInvalid
		return
	}

	role, ok := claims["role"].(string)
	if !ok || role == "" {
		err = errNoRoleClaim
		return
	} else if role != tokenRole {
		err = errWrongRole
		return
	}

	// ID will be interpreted as float. We get an error
	// if we try to do claims["sub"].(int).
	id, ok := claims["sub"].(float64)
	if !ok || id == 0 {
		err = errNoSubClaim
		return
	}

	admin, ok := claims["admin"].(bool)
	if !ok {
		err = errNoAdminClaim
		return
	}

	user = User{ID: int(id), Admin: admin}
	return
}

// newToken generates a new token for the given user and valid for the given duration.
func newToken(user User, duration time.Duration, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"exp":   time.Now().Add(duration).Unix(),
		"admin": user.Admin,
		"role":  role, // Role of the token, either "access" or "refresh"
	})

	return token.SignedString([]byte(hmacSecret))
}

func isInvalidated(pool *redis.Pool, tokenString string) bool {
	c := pool.Get()
	redisToken, _ := c.Do("GET", tokenString)
	return redisToken != nil
}

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

// keyFunc checks if the token was used with the right signing method and
// returns the secret key necessary to parse a JWT token.
func keyFunc(token *jwt.Token) (interface{}, error) {
	// Check if token was signed with the right signing method
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}
	// If yes return the key
	return []byte(hmacSecret), nil
}
