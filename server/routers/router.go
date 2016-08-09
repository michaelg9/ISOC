package routers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"

	"github.com/michaelg9/ISOC/server/controllers"
	"github.com/michaelg9/ISOC/server/middlewares/authentication"
)

// NewRouter creates router for server
func NewRouter(env controllers.Env) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		var handler http.Handler

		// Get the stored handler function for the route
		switch route.Authentication {
		case sessionAuth:
			middlewareEnv := &authentication.MiddlewareEnv{&env}
			handler = negroni.New(
				negroni.HandlerFunc(middlewareEnv.RequireSessionAuth),
				negroni.Wrap(route.HandlerFunc(env)),
			)
		case tokenAuth:
			middlewareEnv := &authentication.MiddlewareEnv{&env}
			handler = negroni.New(
				negroni.HandlerFunc(middlewareEnv.RequireTokenAuth),
				negroni.Wrap(route.HandlerFunc(env)),
			)
		default:
			handler = route.HandlerFunc(env)
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
