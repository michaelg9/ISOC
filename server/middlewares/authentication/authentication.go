package authentication

import (
	"net/http"
	"strings"

	"github.com/michaelg9/ISOC/server/controllers"
)

const (
	errNotAuthorized = "Unauthorized request."

	hmacSecret = "secret"
)

// MiddlewareEnv embeds the controllers.Env struct so that we can write functions on it.
type MiddlewareEnv struct {
	*controllers.Env
}

// RequireSessionAuth is the middleware for routes that require a session to be set with an email.
// TODO: Check if refresh token is valid
func (env *MiddlewareEnv) RequireSessionAuth(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// Set header values to prevent caching
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
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	next(w, r)
}

// RequireTokenAuth is the middleware for routes that require JWT authentication.
func (env *MiddlewareEnv) RequireTokenAuth(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// Set header values to prevent caching
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, errNotAuthorized, http.StatusBadRequest)
		return
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		http.Error(w, errNotAuthorized, http.StatusBadRequest)
		return
	}

	tokenString := authHeaderParts[1]
	user, err := env.Tokens.CheckAccessToken(tokenString)
	if err != nil {
		http.Error(w, errNotAuthorized, http.StatusForbidden)
		return
	}

	// Check if user is registered in database
	_, err = env.DB.GetUser(user)
	if err != nil {
		http.Error(w, errNotAuthorized, http.StatusForbidden)
		return
	}

	next(w, r)
}
