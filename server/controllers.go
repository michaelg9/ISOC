package main

import (
	"fmt"
	"net/http"
)

// Index handles /
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

// Login handles /app/0.1/login?uname={}&pwd{}
func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Login")
}

// Logout handles /app/0.1/logout
func Logout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout")
}

// Upload handles /app/0.1/upload
func Upload(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Upload")
}

// Download handles /data/0.1/q?appid={}
func Download(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Download")
}
