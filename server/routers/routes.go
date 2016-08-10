package routers

import (
	"net/http"

	"github.com/michaelg9/ISOC/server/controllers"
)

const (
	sessionAuth = "Session"
	tokenAuth   = "Token"
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
		func(env controllers.Env) http.HandlerFunc { return env.LoginPage },
	},
	Route{
		"Dashboard",
		"GET",
		"/dashboard",
		sessionAuth,
		func(env controllers.Env) http.HandlerFunc { return env.Dashboard },
	},
	Route{
		"SessionLogin",
		"POST",
		"/login",
		"",
		func(env controllers.Env) http.HandlerFunc { return env.SessionLogin },
	},
	Route{
		"SessionLogout",
		"POST",
		"/logout",
		"",
		func(env controllers.Env) http.HandlerFunc { return env.SessionLogout },
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
		tokenAuth,
		func(env controllers.Env) http.HandlerFunc { return env.Upload },
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
	Route{
		"TokenLogin",
		"POST",
		"/auth/0.1/login",
		"",
		func(env controllers.Env) http.HandlerFunc { return env.TokenLogin },
	},
	Route{
		"Token",
		"POST",
		"/auth/0.1/token",
		"",
		func(env controllers.Env) http.HandlerFunc { return env.Token },
	},
	Route{
		"RefreshToken",
		"POST",
		"/auth/0.1/refresh",
		"",
		func(env controllers.Env) http.HandlerFunc { return env.RefreshToken },
	},
	Route{
		"TokenLogout",
		"POST",
		"/auth/0.1/logout",
		"",
		func(env controllers.Env) http.HandlerFunc { return env.TokenLogout },
	},
	Route{
		"UserData",
		"GET",
		"/data/{email}",
		tokenAuth,
		func(env controllers.Env) http.HandlerFunc { return env.User },
	},
	Route{
		"Device",
		"GET",
		"/data/{email}/{device}",
		tokenAuth,
		func(env controllers.Env) http.HandlerFunc { return env.Device },
	},
	Route{
		"Feature",
		"GET",
		"/data/{email}/{device}/{feature}",
		tokenAuth,
		func(env controllers.Env) http.HandlerFunc { return env.Feature },
	},
}
