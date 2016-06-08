package controllers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"net/http"

	"github.com/michaelg9/ISOC/server/core/mysql"
	"github.com/michaelg9/ISOC/server/services/models"

	"golang.org/x/crypto/bcrypt"
)

// Index handles /
func Index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, "")
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
		// TODO: Return more meaningful error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: Return XML or JSON
	fmt.Fprintf(w, "Success")
}

// LoginWeb renders the template for the log in form of the website
func LoginWeb(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, "")
}

// Logout handles /app/0.1/logout
func Logout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout")
}

// Dashboard handles /dashboard
func Dashboard(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/dashboard.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, "")
}

// Upload handles /app/0.1/upload
func Upload(w http.ResponseWriter, r *http.Request) {
	decoder := xml.NewDecoder(r.Body)

	var d models.DataIn
	if err := decoder.Decode(&d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: more careful parsing
	deviceID := d.Meta.Device
	for _, battery := range d.Battery {
		err := mysql.InsertBatteryData(deviceID, battery.Value, battery.Time)
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

	deviceData, err := mysql.GetDevices(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert DeviceOut to Device
	devices := make([]models.Device, len(deviceData))
	for i, d := range deviceData {
		// TODO: Find a way to make this less ugly
		devices[i].ID = d.ID
		devices[i].Manufacturer = d.Manufacturer
		devices[i].Model = d.Model
		devices[i].OS = d.OS
	}

	// For each device get all its data and append it to the device
	for i, d := range devices {
		var battery []models.Battery
		battery, err = mysql.GetBattery(d.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		devices[i].Battery = battery // TODO: find a way to automatically append all data to device
	}

	// Create the struct for the output data
	response := &models.DataOut{
		Device: devices,
	}

	// Format the output data to JSON
	jsonOut, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(jsonOut))
}
