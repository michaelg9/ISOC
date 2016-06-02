package controllers

import (
	"encoding/json"
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
	// TODO: Check if right parameters entered

	hashedPassword, err := mysql.GetHashedPassword(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if given password fits with stored hash inside the server
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Username: %s Password: %s Success: true", username, password)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: more careful parsing
	deviceID := d.Meta.Device
	for _, battery := range d.Battery {
		err = mysql.InsertBatteryData(deviceID, battery.Value, battery.Time)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	fmt.Fprintln(w, "Success")
}

// Download handles /data/0.1/q
func Download(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("appid")

	dataOut, err := mysql.GetBatteryData(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := &models.Data{
		Battery: *dataOut,
	}

	jsonOut, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(jsonOut))
}
