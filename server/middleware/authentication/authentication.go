package authentication

// TODO: Move session management here

import (
	"net/http"

	"github.com/michaelg9/ISOC/server/core/mysql"
	"github.com/michaelg9/ISOC/server/services/models"
	"golang.org/x/crypto/bcrypt"
)

// RequireBasicAuth is the middleware for routes that require HTTP Basic authentication
func RequireBasicAuth(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	email, password, ok := r.BasicAuth()
	if !ok || email == "" || password == "" {
		http.Error(w, "Unauthorized request", http.StatusUnauthorized)
		return
	}

	// Check if user is registered in database
	var user []models.User
	err := mysql.Get(&user, email, "")
	if err != nil || len(user) == 0 {
		http.Error(w, "Unauthorized request", http.StatusUnauthorized)
		return
	}

	// Check if given password fits with stored hash inside the server
	err = bcrypt.CompareHashAndPassword([]byte(user[0].PasswordHash), []byte(password))
	if err != nil {
		http.Error(w, "Unauthorized request", http.StatusUnauthorized)
		return
	}

	// User is authorised and can proceed
	next(w, r)
}
