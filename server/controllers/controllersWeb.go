package controllers

import (
	"html/template"
	"net/http"
)

// Index handles /
func (env *Env) Index(w http.ResponseWriter, r *http.Request) {
	if err := display(w, "views/index.html", ""); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// LoginWeb handles /login
func (env *Env) LoginWeb(w http.ResponseWriter, r *http.Request) {
	if err := display(w, "views/login.html", ""); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Dashboard handles /dashboard
func (env *Env) Dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store, no-cache, private, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "-1")

	// Get the current user session
	session, err := env.SessionStore.Get(r, "log-in")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the email is set
	email, found := session.Values["email"]
	// If email not set redirect to login page
	if !found || email == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// If email is set go to dashboard
	if err = display(w, "views/dashboard.html", ""); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
