package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/michaelg9/ISOC/server/models"
	"golang.org/x/crypto/bcrypt"
)

// Index handles /
func (env *Env) Index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "views/index.html")
}

// LoginPage handles GET /login
func (env *Env) LoginPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "views/login.html")
}

// SessionLogin handles POST /login
func (env *Env) SessionLogin(w http.ResponseWriter, r *http.Request) {
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
	if err == bcrypt.ErrMismatchedHashAndPassword {
		http.Error(w, errWrongPasswordEmail, http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set session for given user
	session, err := env.SessionStore.Get(r, "log-in")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	refreshToken, err := env.Tokens.NewToken(user, refreshTolkenDelta)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["refreshToken"] = refreshToken
	session.Values["email"] = email
	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	accessToken, err := env.Tokens.NewToken(user, accessTokenDelta)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeResponse(w, "json", models.Tokens{AccessToken: accessToken})
}

// SessionLogout handles /logout
func (env *Env) SessionLogout(w http.ResponseWriter, r *http.Request) {
	// Get the current log-in session of the user
	session, err := env.SessionStore.Get(r, "log-in")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// NOTE: Is this checking necessary?
	// Check if the email is set
	email, found := session.Values["email"]
	// If email not set redirect to login page
	if !found || email == "" {
		http.Error(w, errNoSessionSet, http.StatusUnauthorized)
		return
	}

	token, found := session.Values["refreshToken"]
	tokenString, ok := token.(string)
	if !found || tokenString == "" || !ok {
		http.Error(w, errNoToken, http.StatusInternalServerError)
		return
	}

	if err = env.Tokens.InvalidateToken(tokenString); err != nil {
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

// Dashboard handles /dashboard
func (env *Env) Dashboard(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "views/dashboard.html")
}
