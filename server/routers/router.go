package routers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"

	"github.com/michaelg9/ISOC/server/core/authentication"
)

// NewRouter creates router for server
// TODO: Parameter Env
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		var handler http.Handler

		// Get the stored handler function for the route
		switch route.Authentication {
		case "Basic":
			handler = negroni.New(
				negroni.HandlerFunc(authentication.RequireBasicAuth),
				negroni.Wrap(route.HandlerFunc),
			)
		default:
			handler = route.HandlerFunc
		}

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
