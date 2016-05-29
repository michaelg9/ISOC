package main

import (
	"fmt"
	"net/http"
)

// Index handles /
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

// Login handles /auth/0.1/login
func Login(w http.ResponseWriter, r *http.Request) {
	uname := r.FormValue("username")
	pwd := r.FormValue("password")
	// TODO: Check if non-empty

	fmt.Fprintf(w, "Username: %s Password: %s", uname, pwd)
}

// Logout handles /app/0.1/logout
func Logout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout")
}

// Upload handles /app/0.1/upload
func Upload(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Upload")
}

// Download handles /data/0.1/q
func Download(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("appid")
	// TODO: Check if non-empty
	// TODO: Change query params?

	fmt.Fprintf(w, "API key: %v", key)
}
