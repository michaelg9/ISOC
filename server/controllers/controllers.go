package controllers

// TODO: Change routing with v1..

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/michaelg9/ISOC/server/models"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

const (
	errMissingPasswordOrEmail = "No email and/or password specified."
	errWrongPasswordEmail     = "Wrong password/email combination."
	errNoDeviceID             = "No device ID specified."
	errNoAPIKey               = "No API key specified."
	errNoSessionSet           = "No session set to log-out."
	errMissingToken           = "No token given."

	hmacSecret = "secret"

	accessTokenDelta   = time.Minute * time.Duration(10) // Ten minutes
	refreshTolkenDelta = time.Hour * time.Duration(168)  // One week
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

	// TODO: Move into own file
	response := struct {
		models.DataOut
		User models.User `json:"user"`
	}{
		models.DataOut{Device: devices},
		user,
	}

	err = writeResponse(w, "json", response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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

	err = writeResponse(w, r.FormValue("out"), response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// UpdateUser handles /update/user
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

// Login handles /auth/0.1/login. It returns a short-lived acces token and a long-lived refresh token.
// With the refresh token it is possible to get a new access token via /auth/0.1/token.
func (env *Env) Login(w http.ResponseWriter, r *http.Request) {
	// Get the parameter values for email and password from the URI
	email := r.FormValue("email")
	password := r.FormValue("password")
	// Check if parameters are non-empty
	if email == "" || password == "" {
		http.Error(w, errMissingPasswordOrEmail, http.StatusBadRequest)
		return
	}

	// Get the userdata from the specified email
	user, err := env.DB.GetUser(models.User{Email: email})
	if err == sql.ErrNoRows {
		http.Error(w, errWrongPasswordEmail, http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if given password fits with stored hash inside the server
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		http.Error(w, errWrongPasswordEmail, http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create short-lived acces token
	accessToken, err := env.Tokens.NewToken(user, accessTokenDelta)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create long-lived refresh token
	refreshToken, err := env.Tokens.NewToken(user, refreshTolkenDelta)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	err = writeResponse(w, "json", response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Token handles /auth/0.1/token. It returns a new access token if the user provided
// a valid refresh token.
func (env *Env) Token(w http.ResponseWriter, r *http.Request) {
	// Get the refresh token from the URI
	refreshToken := r.FormValue("refreshToken")
	if refreshToken == "" {
		http.Error(w, errMissingToken, http.StatusBadRequest)
		return
	}

	// Check the validity of the token
	email, err := env.Tokens.CheckToken(refreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// TODO: Check if email is actually saved
	// Create a new access token for the given user
	accessToken, err := env.Tokens.NewToken(models.User{Email: email}, accessTokenDelta)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = writeResponse(w, "json", models.Tokens{AccessToken: accessToken})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Refresh handles /auth/0.1/refresh. The user can get a new refresh token in exchange
// for a valid one. Before the new token gets issued the old one is blacklisted.
func (env *Env) Refresh(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.FormValue("refreshToken")
	if refreshToken == "" {
		http.Error(w, errMissingToken, http.StatusBadRequest)
		return
	}

	email, err := env.Tokens.CheckToken(refreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = env.Tokens.InvalidateToken(refreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newRefreshToken, err := env.Tokens.NewToken(models.User{Email: email}, refreshTolkenDelta)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = writeResponse(w, "json", models.Tokens{RefreshToken: newRefreshToken})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// LogoutToken handles /auth/0.1/logout. The given token will be blacklisted.
func (env *Env) LogoutToken(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	if token == "" {
		http.Error(w, errMissingToken, http.StatusBadRequest)
		return
	}

	err := env.Tokens.InvalidateToken(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "Success")
}

// writeResponse prints a given struct in the specified format to the response writer. If no
// format is specified it will be JSON.
func writeResponse(w http.ResponseWriter, format string, response interface{}) (err error) {
	var out []byte
	switch format {
	case "", "json":
		out, err = json.Marshal(response)
	case "xml":
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
