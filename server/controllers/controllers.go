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
	t, err := template.ParseFiles("views/index.html")
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
	if username == "" || password == "" {
		http.Error(w, "No username and/or password specified.", http.StatusInternalServerError)
		return
	}

	var user []models.User
	err := mysql.Get(&user, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if len(user) == 0 {
		http.Error(w, "Wrong username/password combination.", http.StatusInternalServerError)
		return
	}

	// Check if given password fits with stored hash inside the server
	err = bcrypt.CompareHashAndPassword([]byte(user[0].PasswordHash), []byte(password))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: Return XML or JSON
	fmt.Fprintf(w, "Success")
}

// LoginWeb renders the template for the log in form of the website
func LoginWeb(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("views/login.html")
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
	t, err := template.ParseFiles("views/dashboard.html")
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

	var devicesInfo []models.DeviceStored
	err := mysql.Get(&devicesInfo, key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert DeviceOut to Device
	devices := make([]models.Device, len(devicesInfo))
	for i, d := range devicesInfo {
		devices[i].SetDeviceInfo(d)
	}

	// For each device get all its data and append it to the device
	for i, d := range devices {
		// Get pointers to the arrays which store the tracked data
		// and fill them with the data from the DB
		for _, data := range devices[i].Data.GetContents() {
			err = mysql.Get(data, d.ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
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
