package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

// Index handles /
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

/* TODO: Make login and logout hidden functions and handle
 *       /auth/0.1/{query} in one function
 */

// Login handles /auth/0.1/login/{loginQuery}
func Login(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	q, err := url.ParseQuery(vars["loginQuery"])
	if err != nil {
		// TODO: Proper error handling
		return
	}

	// TODO: Check if right parameters are entered
	fmt.Fprintf(w, "Username: %s Password: %s", q.Get("username"), q.Get("password"))
}

// Logout handles /app/0.1/logout
func Logout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout")
}

// Upload handles /app/0.1/upload
func Upload(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Upload")
}

// Download handles /data/0.1/{query}
func Download(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	q, err := url.ParseQuery(vars["query"])
	if err != nil {
		// TODO: Proper error handling
		return
	}

	fmt.Fprintf(w, "API key: %v", q.Get("appid"))
}
