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
// TODO: Make function behave differently when it's a POST request
func (env *Env) LoginWeb(w http.ResponseWriter, r *http.Request) {
	if err := display(w, "views/login.html", ""); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Dashboard handles /dashboard
func (env *Env) Dashboard(w http.ResponseWriter, r *http.Request) {
	if err := display(w, "views/dashboard.html", ""); err != nil {
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
