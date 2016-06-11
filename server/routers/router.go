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

		handler = route.HandlerFunc

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	// Serve the static contents
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("."))))

	return router
}
