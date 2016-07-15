package authentication

import (
	"net/http"

	"github.com/michaelg9/ISOC/server/controllers"
	"github.com/michaelg9/ISOC/server/services/models"
	"golang.org/x/crypto/bcrypt"
)

const (
	errNotAuthorized = "Unauthorized request."
)

// MiddlewareEnv embeds the controllers.Env struct so that we can write functions on it.
type MiddlewareEnv struct {
	*controllers.Env
}

// RequireBasicAuth is the middleware for routes that require HTTP Basic authentication.
func (env *MiddlewareEnv) RequireBasicAuth(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	email, password, ok := r.BasicAuth()
	if !ok || email == "" || password == "" {
		http.Error(w, errNotAuthorized, http.StatusUnauthorized)
		return
	}

	// Check if user is registered in database
	user, err := env.DB.GetUser(models.User{Email: email})
	if err != nil {
		http.Error(w, errNotAuthorized, http.StatusForbidden)
		return
	}

	// Check if given password fits with stored hash inside the server
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		http.Error(w, errNotAuthorized, http.StatusForbidden)
		return
	}

	// User is authorised and can proceed
	next(w, r)
}

// RequireSessionAuth is the middleware for routes that require a session to be set with an email.
func (env *MiddlewareEnv) RequireSessionAuth(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Set("Cache-Control", "no-store, no-cache, private, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "-1")

	// Get the current user session
	session, err := env.SessionStore.Get(r, "log-in")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the email is set
	email, found := session.Values["email"]
	// If email not set redirect to login page
	if !found || email == "" {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}

	next(w, r)
}
