package routers

import (
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter creates router for server
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		var handler http.Handler

		// Get the stored handler function for the route
		handler = route.HandlerFunc

		router.
			Methods(route.Method). // Define HTTP Method
			Path(route.Pattern).   // Define URI which the function has to handle
			Name(route.Name).      // Define name for this route
			Handler(handler)       // Define handler function
	}

	// Serve the static contents
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("."))))

	return router
}
