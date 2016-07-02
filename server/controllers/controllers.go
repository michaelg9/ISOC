package controllers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"

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
		Path:     "/",
		Secure:   false, // NOTE: Change to true if on actual server
		HttpOnly: true,
	}
}

// Login handles /auth/0.1/login
func Login(w http.ResponseWriter, r *http.Request) {
	// Get the parameter values for email and password from the URI
	email := r.FormValue("email")
	password := r.FormValue("password")
	// Check if parameters are non-empty
	if email == "" || password == "" {
		http.Error(w, "No email and/or password specified.", http.StatusInternalServerError)
		return
	}

	// Get the userdata from the specified email
	var user []models.User
	err := mysql.Get(&user, email, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if len(user) == 0 {
		http.Error(w, "Wrong email/password combination.", http.StatusInternalServerError)
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

	session.Values["email"] = email
	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Success")
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

// SignUp handles /signup
func SignUp(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	if email == "" || password == "" {
		http.Error(w, "You have to specify a password and an email.", http.StatusInternalServerError)
		return
	}

	// Check if the user is already in the database
	var user []models.User
	err := mysql.Get(&user, email, "")
	switch {
	case len(user) == 0:
		// Create new user with hashed password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Insert the credentials of the new user into the database
		user = append(user, models.User{Email: email, PasswordHash: string(hashedPassword)})
		err = mysql.InsertData(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Success")
	case err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	default:
		fmt.Fprintf(w, "User already exists")
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
		err := mysql.InsertData(data, deviceID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	fmt.Fprintln(w, "Success")
}

// InternalDownload handles /data/0.1/user and is only accessible when a user is logged in
func InternalDownload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store, no-cache, private, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "-1")

	// Get the current user session
	session, err := sessionStore.Get(r, "log-in")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the email is set
	email, found := session.Values["email"]
	if !found || email == "" {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	devices, err := getUserDevices(email.(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(devices)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(out))
}

// Download handles /data/0.1/q
func Download(w http.ResponseWriter, r *http.Request) {
	// Get the value of the parameter appid in the URI
	key := r.FormValue("appid")
	if key == "" {
		http.Error(w, "No API Key given.", http.StatusInternalServerError)
	}

	var user []models.User
	err := mysql.Get(&user, "", key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := getUserDevices(user[0].Email)

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

	fmt.Fprint(w, string(out))
}

func getUserDevices(email string) (models.DataOut, error) {
	// Query the database for the data which belongs to the API key
	var devicesInfo []models.DeviceStored
	err := mysql.Get(&devicesInfo, email)
	if err != nil {
		return models.DataOut{}, err
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
				return models.DataOut{}, err
			}
		}
	}

	result := models.DataOut{
		Device: devices,
	}

	return result, nil
}
