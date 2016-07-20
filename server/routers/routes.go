package routers

import (
	"net/http"

	"github.com/michaelg9/ISOC/server/controllers"
)

const (
	basicAuth   = "Basic"
	sessionAuth = "Session"
)

// Route struct for creating routes
type Route struct {
	Name           string
	Method         string                                 // HTTP method
	Pattern        string                                 // URI for this route
	Authentication string                                 // Authentication method to be used
	HandlerFunc    func(controllers.Env) http.HandlerFunc // Handler function specified in controllers
}

// Routes is an array of routes
type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		"",
		func(env controllers.Env) http.HandlerFunc { return env.Index },
	},
	Route{
		"Login",
		"GET",
		"/login",
		"",
		func(env controllers.Env) http.HandlerFunc { return env.LoginGET },
	},
	Route{
		"Dashboard",
		"GET",
		"/dashboard",
		sessionAuth,
		func(env controllers.Env) http.HandlerFunc { return env.Dashboard },
	},
	Route{
		"Login",
		"POST",
		"/login",
		"",
		func(env controllers.Env) http.HandlerFunc { return env.LoginPOST },
	},
	Route{
		"Logout",
		"POST",
		"/logout",
		"",
		func(env controllers.Env) http.HandlerFunc { return env.Logout },
	},
	Route{
		"SignUp",
		"POST",
		"/signup",
		"",
		func(env controllers.Env) http.HandlerFunc { return env.SignUp },
	},
	Route{
		"Upload",
		"POST",
		"/app/0.1/upload",
		basicAuth,
		func(env controllers.Env) http.HandlerFunc { return env.Upload },
	},
	Route{
		"Download",
		"GET",
		"/data/0.1/q",
		"",
		func(env controllers.Env) http.HandlerFunc { return env.Download },
	},
	Route{
		"InternalDownload",
		"GET",
		"/data/0.1/user",
		sessionAuth,
		func(env controllers.Env) http.HandlerFunc { return env.InternalDownload },
	},
	Route{
		"UpdateUser",
		"POST",
		"/update/user",
		sessionAuth,
		func(env controllers.Env) http.HandlerFunc { return env.UpdateUser },
	},
}
