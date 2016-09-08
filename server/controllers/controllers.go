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
	"strings"

	"github.com/michaelg9/ISOC/server/models"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

const (
	errNoPasswordOrEmail   = "No email and/or password specified."
	errWrongPasswordEmail  = "Wrong password/email combination."
	errNoDeviceID          = "No device ID specified."
	errNoSessionSet        = "No session set to log-out."
	errNoToken             = "No token given."
	errTokenInvalid        = "Token is invalid."
	errTokenAlreadyInvalid = "Token already invalid."
	errWrongUser           = "No such user."
	errDeviceIDNotInt      = "The device value is not an int."
	errUserIDNotInt        = "The user value is not an int."
	errWrongFeature        = "This device has no data for this feature."
	errWrongDeviceOrUser   = "The given user ID/device ID combination is not valid."
	errForbidden           = "Forbidden resource access."
	errUserExists          = "User already exists."
	errNoDevice            = "No device data specified."

	hmacSecret = "secret"

	// 0 is the user ID when no user is specified
	noUser = 0
)

// key type for use with context
type key int

// UserKey is the key to the user data in the context
const UserKey key = 0

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
	_, err := env.DB.GetUser(models.User{Email: email})
	if err == nil {
		// User already exists
		http.Error(w, errUserExists, http.StatusConflict)
		return
	} else if err != sql.ErrNoRows {
		// Unexpected error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	decoder := xml.NewDecoder(r.Body)

	// Decode the given XML from the request body into the struct defined in models
	var d models.Upload
	if err = decoder.Decode(&d); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create new user with hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert the credentials of the new user into the database
	userID, err := env.DB.CreateUser(models.User{Email: email, PasswordHash: string(hashedPassword)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	device := d.Meta
	noDevice := device == models.AboutDevice{}
	if noDevice {
		http.Error(w, errNoDevice, http.StatusBadRequest)
		return
	}

	deviceID, err := env.DB.CreateDeviceForUser(models.User{ID: userID}, device)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, deviceID)
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

	deviceID := d.Meta.ID
	// Check if the device ID was specified in the input.
	// If no ID was specified it defaults to 0.
	if deviceID == 0 {
		http.Error(w, errNoDeviceID, http.StatusBadRequest)
		return
	}

	// If decoding was successfull input the data into the database
	for _, data := range d.Features.GetContents() {
		err := env.DB.CreateFeatureForDevice(models.AboutDevice{ID: deviceID}, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	fmt.Fprint(w, http.StatusText(http.StatusOK))
}

// UpdateUser handles /update/user
func (env *Env) UpdateUser(w http.ResponseWriter, r *http.Request) {
	user, ok := context.Get(r, UserKey).(models.User)
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Overwrite the stored values with the values given in the request.
	// If a value isn't given it defaults to the empty string, thus we don't update it.
	user.Email = r.FormValue("email")
	user.PasswordHash = r.FormValue("password")

	if user.PasswordHash != "" {
		// If password is non-empty create a new hash for the database.
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		user.PasswordHash = string(hashedPassword)
	}

	if err := env.DB.UpdateUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, http.StatusText(http.StatusOK))
}

// GetUser handles /data/{user}. It gets all the devices and its data from
// the user and information about the user.
func (env *Env) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["user"])
	if err != nil {
		http.Error(w, errUserIDNotInt, http.StatusInternalServerError)
		return
	}

	if !userIsAuthorized(r, userID) {
		http.Error(w, errForbidden, http.StatusForbidden)
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

// GetAllDevices handels /data/all/devices. It gets the info about every stored device.
func (env *Env) GetAllDevices(w http.ResponseWriter, r *http.Request) {
	isAdmin := userIsAuthorized(r, noUser)
	if !isAdmin {
		http.Error(w, errForbidden, http.StatusForbidden)
		return
	}

	response, err := env.DB.GetAllDevices()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = writeResponse(w, r.FormValue("out"), response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetDevice handles /data/{user}/{device}. It gets all the info about the
// specified device of the given user. Device should be the device ID (for now).
func (env *Env) GetDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["user"])
	if err != nil {
		http.Error(w, errUserIDNotInt, http.StatusInternalServerError)
		return
	}

	if !userIsAuthorized(r, userID) {
		http.Error(w, errForbidden, http.StatusForbidden)
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

// GetAllFeatures handles /data/all/features. It gets every entry from all features.
func (env *Env) GetAllFeatures(w http.ResponseWriter, r *http.Request) {
	isAdmin := userIsAuthorized(r, noUser)
	if !isAdmin {
		http.Error(w, errForbidden, http.StatusForbidden)
		return
	}

	var response models.Features
	for _, feature := range response.GetContents() {
		err := env.DB.GetAllFeatureData(feature)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err := writeResponse(w, r.FormValue("out"), response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetAllOfFeature handles /data/all/features/{feature}. It gets every entry from the given
// feature.
func (env *Env) GetAllOfFeature(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Capitalize first letter of feature name
	feature := strings.Title(vars["feature"])

	isAdmin := userIsAuthorized(r, noUser)
	if !isAdmin {
		http.Error(w, errForbidden, http.StatusForbidden)
		return
	}

	// Get all data of the given feature
	response, err := getFeatureResponse(env, feature)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = writeResponse(w, r.FormValue("out"), response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetFeature handles /data/{user}/{device}/{feature}. It gets all the saved data from the
// specified feature from the given device of the user.
func (env *Env) GetFeature(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["user"])
	if err != nil {
		http.Error(w, errUserIDNotInt, http.StatusInternalServerError)
		return
	}

	if !userIsAuthorized(r, userID) {
		http.Error(w, errForbidden, http.StatusForbidden)
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

	// Capitalize first letter of feature name
	feature := strings.Title(vars["feature"])
	response, err := getFeatureResponse(env, feature, deviceID)
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

// getFeatureResponse queries the database to get the feature with the given name from the device
// with the specified device ID.
func getFeatureResponse(env *Env, feature string, args ...int) (response models.Features, err error) {
	// Get the reflect value of the field with the name of the feature
	featureValue := reflect.ValueOf(&response).Elem().FieldByName(feature)
	// Get a pointer to the feature value
	featurePtr := featureValue.Addr().Interface()

	// If a device ID was specified load only the data that is associated with the device,
	// otherwise load all the data from the given feature.
	if len(args) == 1 {
		// Get the feature data from the specified device and save it to the response struct
		deviceID := args[0]
		err = env.DB.GetFeatureOfDevice(models.AboutDevice{ID: deviceID}, featurePtr)
	} else {
		err = env.DB.GetAllFeatureData(featurePtr)
	}
	return
}

// userIsAuthorized checks if the given user is allowed access to the resource i.e. when the user ID
// in the request matches the ID from the logged in user or when the logged in user is an admin.
func userIsAuthorized(r *http.Request, givenID int) bool {
	user, ok := context.Get(r, UserKey).(models.User)
	if !ok {
		return false
	}

	isAdmin := user.Admin
	rightID := givenID == user.ID
	return rightID || isAdmin
}

// writeResponse prints a given struct in the specified format to the response writer. If no
// format is specified it will be JSON.
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
