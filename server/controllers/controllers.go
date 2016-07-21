package controllers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/michaelg9/ISOC/server/services/models"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

const (
	errMissingPasswordOrEmail = "No email and/or password specified."
	errWrongPasswordEmail     = "Wrong password/email combination."
	errNoDeviceID             = "No device ID specified."
	errNoAPIKey               = "No API key specified."
	errNoSessionSet           = "No session set to log-out."
)

// Env contains the environment information
type Env struct {
	DB           models.Datastore
	SessionStore *sessions.CookieStore
}

// SignUp handles /signup
func (env *Env) SignUp(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	if email == "" || password == "" {
		http.Error(w, errMissingPasswordOrEmail, http.StatusBadRequest)
		return
	}

	// Check if the user is already in the database
	user, err := env.DB.GetUser(models.User{Email: email})
	switch {
	// User does not exist
	case user == models.User{}:
		// Create new user with hashed password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Insert the credentials of the new user into the database
		err = env.DB.CreateUser(models.User{Email: email, PasswordHash: string(hashedPassword)})
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
func (env *Env) Upload(w http.ResponseWriter, r *http.Request) {
	decoder := xml.NewDecoder(r.Body)

	// Decode the given XML from the request body into the struct defined in models
	var d models.DataIn
	if err := decoder.Decode(&d); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	deviceID := d.Meta.Device
	// Check if the device ID was specified in the input.
	// If no ID was specified it defaults to 0.
	if deviceID == 0 {
		http.Error(w, errNoDeviceID, http.StatusBadRequest)
		return
	}

	// If decoding was successfull input the data into the database
	for _, data := range d.DeviceData.GetContents() {
		err := env.DB.CreateData(models.DeviceStored{ID: deviceID}, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	fmt.Fprint(w, "Success")
}

// InternalDownload handles /data/0.1/user and is only accessible when a user is logged in
func (env *Env) InternalDownload(w http.ResponseWriter, r *http.Request) {
	// Because of the middleware we know that these values exist
	session, _ := env.SessionStore.Get(r, "log-in")
	email := session.Values["email"]
	devices, err := env.DB.GetDevicesFromUser(models.User{Email: email.(string)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user, err := env.DB.GetUser(models.User{Email: email.(string)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Make sure we don't send the password hash over the wire
	user.PasswordHash = ""

	outStruct := struct {
		models.DataOut
		User models.User `json:"user"`
	}{
		models.DataOut{Device: devices},
		user,
	}

	out, err := json.Marshal(outStruct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(out))
}

// Download handles /data/0.1/q
func (env *Env) Download(w http.ResponseWriter, r *http.Request) {
	// Get the value of the parameter appid in the URI
	key := r.FormValue("appid")
	if key == "" {
		http.Error(w, errNoAPIKey, http.StatusBadRequest)
		return
	}

	devices, err := env.DB.GetDevicesFromUser(models.User{APIKey: key})
	response := models.DataOut{
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

	fmt.Fprint(w, string(out))
}

// UpdateUser handles /update/user
func (env *Env) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Because of the middleware we know that these values exist
	session, _ := env.SessionStore.Get(r, "log-in")
	oldEmail := session.Values["email"]

	newEmail := r.FormValue("email")
	password := r.FormValue("password")

	user, err := env.DB.GetUser(models.User{Email: oldEmail.(string)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if newEmail != "" {
		user.Email = newEmail
		// TODO: Don't use magic strings
		err = env.DB.UpdateUser(user, "Email")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Update the current session with the new email address
		session.Values["email"] = newEmail
		if err = session.Save(r, w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user.PasswordHash = string(hashedPassword)
		err = env.DB.UpdateUser(user, "PasswordHash")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	fmt.Fprint(w, "Success")
}
