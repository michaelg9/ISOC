package controllers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"github.com/michaelg9/ISOC/server/services/models"
	"golang.org/x/crypto/bcrypt"
)

// Index handles /
func (env *Env) Index(w http.ResponseWriter, r *http.Request) {
	if err := display(w, "views/index.html", ""); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// LoginGET handles GET /login
func (env *Env) LoginGET(w http.ResponseWriter, r *http.Request) {
	if err := display(w, "views/login.html", ""); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// LoginPOST handles POST /login
func (env *Env) LoginPOST(w http.ResponseWriter, r *http.Request) {
	// TODO: Refactor into function "checkUser"
	// Get the parameter values for email and password from the URI
	email := r.FormValue("email")
	password := r.FormValue("password")
	// Check if parameters are non-empty
	if email == "" || password == "" {
		http.Error(w, errMissingPasswordOrEmail, http.StatusBadRequest)
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

	session.Values["email"] = email
	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Success")
}

// Logout handles /logout
func (env *Env) Logout(w http.ResponseWriter, r *http.Request) {
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
	if err := display(w, "views/dashboard.html", ""); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// display takes a filepath to an HTML or template, passes the given
// data to it and displays it
func display(w http.ResponseWriter, filePath string, data interface{}) error {
	t, err := template.ParseFiles(filePath)
	if err != nil {
		return err
	}
	t.Execute(w, data)
	return nil
}
