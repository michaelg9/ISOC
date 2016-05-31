package controllers

import (
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/michaelg9/ISOC/server/core/mysql"
	"github.com/michaelg9/ISOC/server/services/models"
	"golang.org/x/crypto/bcrypt"
)

// Index handles /
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

// Login handles /auth/0.1/login
func Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	// TODO: Check if non-empty

	hashedPassword := mysql.GetHashedPassword(username)
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	var success bool
	if err != nil {
		// TODO: Error handling
		success = false
	} else {
		success = true
	}

	fmt.Fprintf(w, "Username: %s Password: %s Success: %t", username, password, success)
}

// Logout handles /app/0.1/logout
func Logout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout")
}

// Upload handles /app/0.1/upload
func Upload(w http.ResponseWriter, r *http.Request) {
	decoder := xml.NewDecoder(r.Body)

	var d models.Data
	err := decoder.Decode(&d)
	if err != nil {
		//TODO: Error handling
	}

	fmt.Fprintln(w, d)
}

// Download handles /data/0.1/q
func Download(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("appid")
	// TODO: Check if non-empty

	fmt.Fprintf(w, "API key: %v", key)
}
