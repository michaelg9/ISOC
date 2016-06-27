package controllers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"net/http"
	"reflect"

	"github.com/michaelg9/ISOC/server/core/mysql"
	"github.com/michaelg9/ISOC/server/services/models"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var sessionStore *sessions.CookieStore

func init() {
	// Initiate the a cookie store to save our sessions
	sessionStore = sessions.NewCookieStore([]byte("something-very-secret"))
	sessionStore.Options = &sessions.Options{
		Path:     "/dashboard",
		Secure:   false, // NOTE: Change to true if on actual server
		HttpOnly: true,
	}
}

// Index handles /
func Index(w http.ResponseWriter, r *http.Request) {
	if err := display(w, "views/index.html", ""); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Login handles /auth/0.1/login
func Login(w http.ResponseWriter, r *http.Request) {
	// Get the parameter values for username and password from the URI
	username := r.FormValue("username")
	password := r.FormValue("password")
	// Check if parameters are non-empty
	if username == "" || password == "" {
		http.Error(w, "No username and/or password specified.", http.StatusInternalServerError)
		return
	}

	// Get the userdata from the specified username
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

	// Set session for given user
	session, err := sessionStore.Get(r, "log-in")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["username"] = username
	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Success")
}

// LoginWeb handles /login
func LoginWeb(w http.ResponseWriter, r *http.Request) {
	if err := display(w, "views/login.html", ""); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Logout handles /auth/0.1/logout
func Logout(w http.ResponseWriter, r *http.Request) {
	// Get the current log-in session of the user
	session, err := sessionStore.Get(r, "log-in")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set MaxAge to -1 to delete the session
	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Success")
}

// Dashboard handles /dashboard
func Dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store, no-cache, private, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "-1")

	// Get the current user session
	session, err := sessionStore.Get(r, "log-in")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the username is set
	username, found := session.Values["username"]
	// If username not set redirect to login page
	if !found || username == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// If username is set go to dashboard
	if err = display(w, "views/dashboard.html", ""); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Upload handles /app/0.1/upload
func Upload(w http.ResponseWriter, r *http.Request) {
	decoder := xml.NewDecoder(r.Body)

	// Decode the given XML from the request body into the struct defined in models
	var d models.DataIn
	if err := decoder.Decode(&d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	deviceID := d.Meta.Device
	// Check if the device ID was specified in the input.
	// If no ID was specified it defaults to 0.
	if deviceID == 0 {
		http.Error(w, "No device ID specified.", http.StatusInternalServerError)
		return
	}

	// If decoding was successfull input the data into the database
	for _, data := range d.DeviceData.GetContents() {
		// TODO: Check if this is the right place or if you should change
		//       the InsertData function
		dataValue := reflect.Indirect(reflect.ValueOf(data)).Interface()
		err := mysql.InsertData(deviceID, dataValue)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	fmt.Fprintln(w, "Success")
}

// Download handles /data/0.1/q
func Download(w http.ResponseWriter, r *http.Request) {
	// Get the value of the parameter appid in the URI
	key := r.FormValue("appid")

	// Query the database for the data which belongs to the API key
	var devicesInfo []models.DeviceStored
	err := mysql.Get(&devicesInfo, key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert DeviceOut to Device
	devices := make([]models.Device, len(devicesInfo))
	for i, d := range devicesInfo {
		devices[i].DeviceInfo = d
	}

	// For each device get all its data and append it to the device
	for i, d := range devices {
		// Get pointers to the arrays which store the tracked data
		// and fill them with the data from the DB
		for _, data := range devices[i].Data.GetContents() {
			err = mysql.Get(data, d.DeviceInfo.ID)
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

	// Format the output into either JSON or XML according to the
	// value specified from the parameter "out". The default value
	// is JSON.
	var out []byte
	switch r.FormValue("out") {
	case "", "json":
		out, err = json.Marshal(response)
	case "xml":
		out, err = xml.Marshal(response)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(out))
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
