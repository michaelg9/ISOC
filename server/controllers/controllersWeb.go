package controllers

import "net/http"

// Index handles /
func (env *Env) Index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "views/index.html")
}

// LoginPage handles GET /login
func (env *Env) LoginPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "views/login.html")
}

// Dashboard handles /dashboard
func (env *Env) Dashboard(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "views/dashboard.html")
}

// Admin handles /dashboard/admin
func (env *Env) Admin(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "views/admin.html")
}
