package controllers

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/michaelg9/ISOC/server/models"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

const (
	errNoPasswordOrEmail   = "No email and/or password specified."
	errWrongPasswordEmail  = "Wrong password/email combination."
	errNoDeviceID          = "No device ID specified."
	errNoAPIKey            = "No API key specified."
	errWrongAPIKey         = "Wrong API key."
	errNoSessionSet        = "No session set to log-out."
	errNoToken             = "No token given."
	errTokenInvalid        = "Token is invalid."
	errTokenAlreadyInvalid = "Token already invalid."
	errWrongUser           = "No such user."
	errDeviceIDNotInt      = "The device value is not an int."
	errUserIDNotInt        = "The user value is not an int."
	errWrongFeature        = "This device has no data for this feature."
	errWrongDeviceOrUser   = "The given user ID/device ID combination is not valid."

	hmacSecret = "secret"
)

// Env contains the environment information
type Env struct {
	DB           models.Datastore
	Tokens       models.TokenControl
	SessionStore *sessions.CookieStore
}

// SignUp handles /signup
func (env *Env) SignUp(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	if email == "" || password == "" {
		http.Error(w, errNoPasswordOrEmail, http.StatusBadRequest)
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
	var d models.Upload
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
	for _, data := range d.TrackedData.GetContents() {
		err := env.DB.CreateData(models.AboutDevice{ID: deviceID}, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	fmt.Fprint(w, "Success")
}

// UpdateUser handles /update/user
// TODO: Use Token Auth
func (env *Env) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Because of the middleware we know that these values exist
	session, _ := env.SessionStore.Get(r, "log-in")
	oldEmail := session.Values["email"]

	// We need to get the user from the database to get its user ID
	user, err := env.DB.GetUser(models.User{Email: oldEmail.(string)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Overwrite the stored values with the values given in the request.
	// If a value isn't given it defaults to the empty string, thus we don't update it.
	user.Email = r.FormValue("email")
	user.PasswordHash = r.FormValue("password")
	user.APIKey = r.FormValue("apiKey")

	if user.Email != "" {
		// Update the current session with the new email address
		session.Values["email"] = user.Email
		if err = session.Save(r, w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if user.PasswordHash != "" {
		// If password is non-empty create a new hash for the database.
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		user.PasswordHash = string(hashedPassword)
	}

	// For the API key we use 0/1 for false/true since the actual value of the key will be determined
	// by the MySQL database by calling the UUID() function. Hence if the value is not 1 we set the
	// value in the struct to be the empty string so that the apiKey won't be updated.
	if user.APIKey != "1" {
		user.APIKey = ""
	}

	if err = env.DB.UpdateUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "Success")
}

// User handles /data/{user}. It gets all the devices and its data from
// the user and information about the user.
// TODO: Check if userID fits with token
func (env *Env) User(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["user"])
	if err != nil {
		http.Error(w, errUserIDNotInt, http.StatusInternalServerError)
		return
	}

	user, err := env.DB.GetUser(models.User{ID: userID})
	if err == sql.ErrNoRows {
		http.Error(w, errWrongUser, http.StatusInternalServerError)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	devices, err := env.DB.GetDevicesFromUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.UserResponse{
		User:    user,
		Devices: devices,
	}

	err = writeResponse(w, r.FormValue("out"), response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Device handles /data/{user}/{device}. It gets all the info about the
// specified device of the given user. Device should be the device ID (for now).
func (env *Env) Device(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["user"])
	if err != nil {
		http.Error(w, errUserIDNotInt, http.StatusInternalServerError)
		return
	}

	deviceID, err := strconv.Atoi(vars["device"])
	if err != nil {
		http.Error(w, errDeviceIDNotInt, http.StatusInternalServerError)
		return
	}

	device := models.Device{AboutDevice: models.AboutDevice{ID: deviceID}}
	response, err := env.DB.GetDeviceFromUser(models.User{ID: userID}, device)
	if err == sql.ErrNoRows {
		http.Error(w, errWrongDeviceOrUser, http.StatusInternalServerError)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = writeResponse(w, r.FormValue("out"), response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Feature handles /data/{user}/{device}/{feature}. It gets all the saved data from the
// specified feature from the given device of the user.
func (env *Env) Feature(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["user"])
	if err != nil {
		http.Error(w, errUserIDNotInt, http.StatusInternalServerError)
		return
	}

	deviceID, err := strconv.Atoi(vars["device"])
	if err != nil {
		http.Error(w, errDeviceIDNotInt, http.StatusInternalServerError)
		return
	}

	// Check if device ID is registered with users
	deviceStruct := models.Device{AboutDevice: models.AboutDevice{ID: deviceID}}
	_, err = env.DB.GetDeviceFromUser(models.User{ID: userID}, deviceStruct)
	if err == sql.ErrNoRows {
		http.Error(w, errWrongDeviceOrUser, http.StatusInternalServerError)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: Refactor into new function and add comments
	var response models.TrackedData
	v := reflect.ValueOf(&response).Elem()
	f := vars["feature"]
	ptrToFeature := v.FieldByName(f).Addr().Interface()
	err = env.DB.GetData(models.AboutDevice{ID: deviceID}, ptrToFeature)
	if err == sql.ErrNoRows {
		http.Error(w, errWrongFeature, http.StatusInternalServerError)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = writeResponse(w, r.FormValue("out"), response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// writeResponse prints a given struct in the specified format to the response writer. If no
// format is specified it will be JSON.
// TODO: Add csv support
func writeResponse(w http.ResponseWriter, format string, response interface{}) (err error) {
	var out []byte
	switch format {
	case "", "json":
		w.Header().Set("Content-Type", "application/json")
		out, err = json.Marshal(response)
	case "xml":
		w.Header().Set("Content-Type", "application/xml")
		out, err = xml.Marshal(response)
	default:
		err = errors.New("Output type not supported.")
	}
	if err != nil {
		return
	}

	fmt.Fprint(w, string(out))
	return
}
