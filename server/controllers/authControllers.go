package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/michaelg9/ISOC/server/models"
	"golang.org/x/crypto/bcrypt"
)

// TokenLogin handles /auth/0.1/login. It returns a short-lived acces token and a long-lived refresh token.
// With the refresh token it is possible to get a new access token via /auth/0.1/token.
func (env *Env) TokenLogin(w http.ResponseWriter, r *http.Request) {
	// Set header values to prevent caching
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")

	// Get the parameter values for email and password from the URI
	email := r.FormValue("email")
	password := r.FormValue("password")
	// Check if parameters are non-empty
	if email == "" || password == "" {
		http.Error(w, errNoPasswordOrEmail, http.StatusBadRequest)
		return
	}

	// Get the userdata from the specified email
	user, err := env.DB.GetUser(models.User{Email: email})
	if err == sql.ErrNoRows {
		http.Error(w, errWrongPasswordEmail, http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if given password fits with stored hash inside the server
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword || err == bcrypt.ErrHashTooShort {
		http.Error(w, errWrongPasswordEmail, http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create short-lived acces token
	accessToken, err := env.Tokens.NewToken(user, accessTokenDelta)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create long-lived refresh token
	refreshToken, err := env.Tokens.NewToken(user, refreshTolkenDelta)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	err = writeResponse(w, "json", response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Token handles /auth/0.1/token. It returns a new access token if the user provided
// a valid refresh token.
func (env *Env) Token(w http.ResponseWriter, r *http.Request) {
	// Set header values to prevent caching
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")

	var email string
	var err error
	// Get the refresh token from the URI
	refreshToken := r.FormValue("refreshToken")
	if refreshToken == "" {
		// See if there is a valid refresh token in the cookie store
		email, err = refreshTokenInCookie(env, w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//http.Error(w, errNoToken, http.StatusBadRequest)
		//return
	} else {
		// Check the validity of the token
		email, err = env.Tokens.CheckToken(refreshToken)
		if err != nil {
			http.Error(w, errTokenInvalid, http.StatusUnauthorized)
			return
		}
	}

	// Create a new access token for the given user
	accessToken, err := env.Tokens.NewToken(models.User{Email: email}, accessTokenDelta)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = writeResponse(w, "json", models.Tokens{AccessToken: accessToken})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// RefreshToken handles /auth/0.1/refresh. The user can get a new refresh token in exchange
// for a valid one. Before the new token gets issued the old one is blacklisted.
func (env *Env) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// Set header values to prevent caching
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")

	refreshToken := r.FormValue("refreshToken")
	if refreshToken == "" {
		http.Error(w, errNoToken, http.StatusBadRequest)
		return
	}

	email, err := env.Tokens.CheckToken(refreshToken)
	if err != nil {
		http.Error(w, errTokenInvalid, http.StatusUnauthorized)
		return
	}

	err = env.Tokens.InvalidateToken(refreshToken)
	if err != nil {
		http.Error(w, errTokenAlreadyInvalid, http.StatusInternalServerError)
		return
	}

	newRefreshToken, err := env.Tokens.NewToken(models.User{Email: email}, refreshTolkenDelta)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = writeResponse(w, "json", models.Tokens{RefreshToken: newRefreshToken})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// TokenLogout handles /auth/0.1/logout. The given token will be blacklisted.
func (env *Env) TokenLogout(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	if token == "" {
		http.Error(w, errNoToken, http.StatusBadRequest)
		return
	}

	err := env.Tokens.InvalidateToken(token)
	if err != nil {
		http.Error(w, errTokenAlreadyInvalid, http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "Success")
}

// refreshTokenInCookie checks if there is a valid refresh token in
func refreshTokenInCookie(env *Env, w http.ResponseWriter, r *http.Request) (email string, err error) {
	session, err := env.SessionStore.Get(r, "log-in")
	if err != nil {
		return
	}

	refreshToken, found := session.Values["refreshToken"]
	tokenString, ok := refreshToken.(string)
	if !(found && ok) {
		err = errors.New("No refresh token found")
		return
	}

	email, err = env.Tokens.CheckToken(tokenString)
	if err != nil {
		return
	}

	newRefreshToken, err := env.Tokens.NewToken(models.User{Email: email}, refreshTolkenDelta)
	if err != nil {
		return
	}

	// Invalidate old token
	if err = env.Tokens.InvalidateToken(tokenString); err != nil {
		return
	}

	// Update refresh token so the user does not get logged out when the old token would have expired
	session.Values["refreshToken"] = newRefreshToken
	if err = session.Save(r, w); err != nil {
		return
	}

	return
}
