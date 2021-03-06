package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/michaelg9/ISOC/server/models"
	"golang.org/x/crypto/bcrypt"
)

/* Handlers for token authentication */

// TokenLogin handles /auth/0.1/login. It returns a short-lived acces token and a long-lived refresh token.
// With the refresh token it is possible to get a new access token via /auth/0.1/token.
func (env *Env) TokenLogin(w http.ResponseWriter, r *http.Request) {
	// Set header values to prevent caching
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")

	user, code := getUser(env, w, r)
	if code != http.StatusOK {
		http.Error(w, http.StatusText(code), code)
		return
	}

	tokens, err := getTokens(env, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = writeResponse(w, "json", tokens)
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

	var user models.User
	var err error
	// Get the refresh token from the URI
	refreshToken := r.FormValue("refreshToken")
	if refreshToken == "" {
		// See if there is a valid refresh token in the cookie store
		user, err = refreshTokenInCookie(env, w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// Check the validity of the token
		user, err = env.Tokens.CheckRefreshToken(refreshToken)
		if err != nil {
			http.Error(w, errTokenInvalid, http.StatusUnauthorized)
			return
		}
	}

	// Create a new access token for the given user
	accessToken, err := env.Tokens.NewAccessToken(user)
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

	oldRefreshToken := r.FormValue("refreshToken")
	if oldRefreshToken == "" {
		http.Error(w, errNoToken, http.StatusBadRequest)
		return
	}

	_, refreshToken, err := newRefreshToken(env, oldRefreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = writeResponse(w, "json", models.Tokens{RefreshToken: refreshToken})
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

	fmt.Fprint(w, http.StatusText(http.StatusOK))
}

/* Handlers for web session */

// SessionLogin handles POST /login
func (env *Env) SessionLogin(w http.ResponseWriter, r *http.Request) {
	user, code := getUser(env, w, r)
	if code != http.StatusOK {
		http.Error(w, http.StatusText(code), code)
		return
	}

	tokens, err := getTokens(env, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set session for given user
	session, err := env.SessionStore.Get(r, "log-in")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["refreshToken"] = tokens.RefreshToken
	session.Values["id"] = user.ID
	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeResponse(w, "json", models.LoginResponse{AccessToken: tokens.AccessToken, ID: user.ID})
}

// SessionLogout handles /logout
func (env *Env) SessionLogout(w http.ResponseWriter, r *http.Request) {
	// Get the current log-in session of the user
	session, err := env.SessionStore.Get(r, "log-in")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, ok := session.Values["refreshToken"].(string)
	if !ok || token == "" {
		return // We don't throw an error because the user is already logged out
	}

	if err = env.Tokens.InvalidateToken(token); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set MaxAge to -1 to delete the session
	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Success")
}

// getUser gets the given user from the parameters specified in the HTTP request.
func getUser(env *Env, w http.ResponseWriter, r *http.Request) (user models.User, responseCode int) {
	// Get the parameter values for email and password from the URI
	email := r.FormValue("email")
	password := r.FormValue("password")
	// Check if parameters are non-empty
	if email == "" || password == "" {
		return user, http.StatusBadRequest
	}

	// Get the userdata from the specified email
	user, err := env.DB.GetUser(models.User{Email: email})
	if err == sql.ErrNoRows {
		return user, http.StatusUnauthorized
	} else if err != nil {
		return user, http.StatusInternalServerError
	}

	// Check if given password fits with stored hash inside the server
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword || err == bcrypt.ErrHashTooShort {
		return user, http.StatusUnauthorized
	} else if err != nil {
		return user, http.StatusInternalServerError
	}

	return user, http.StatusOK
}

func getTokens(env *Env, user models.User) (tokens models.Tokens, err error) {
	// Create short-lived acces token
	accessToken, err := env.Tokens.NewAccessToken(user)
	if err != nil {
		return
	}

	// Create long-lived refresh token
	refreshToken, err := env.Tokens.NewRefreshToken(user)
	if err != nil {
		return
	}

	tokens = models.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return
}

// refreshTokenInCookie checks if there is a valid refresh token in
func refreshTokenInCookie(env *Env, w http.ResponseWriter, r *http.Request) (user models.User, err error) {
	session, err := env.SessionStore.Get(r, "log-in")
	if err != nil {
		return
	}

	// Get refresh token saved inside cookies
	oldRefreshToken, found := session.Values["refreshToken"]
	oldRefreshTokenString, ok := oldRefreshToken.(string)
	if !(found && ok) {
		err = errors.New("No refresh token found")
		return
	}

	user, refreshToken, err := newRefreshToken(env, oldRefreshTokenString)
	if err != nil {
		return
	}

	// Update refresh token so the user does not get logged out when the old token would have expired
	session.Values["refreshToken"] = refreshToken
	if err = session.Save(r, w); err != nil {
		return
	}

	return
}

// newRefreshToken returns a new refresh token in exchange for a new old one which will be invalidated.
func newRefreshToken(env *Env, oldRefreshToken string) (user models.User, refreshToken string, err error) {
	// Check if old token is valid
	user, err = env.Tokens.CheckRefreshToken(oldRefreshToken)
	if err != nil {
		return
	}

	// Invalidate old token
	if err = env.Tokens.InvalidateToken(oldRefreshToken); err != nil {
		return
	}

	// Create new refresh token
	refreshToken, err = env.Tokens.NewRefreshToken(user)
	if err != nil {
		return
	}

	return
}
