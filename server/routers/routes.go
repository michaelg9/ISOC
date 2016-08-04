package routers

import (
	"net/http"

	"github.com/michaelg9/ISOC/server/controllers"
)

const (
	basicAuth   = "Basic"
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
		"LoginToken",
		"POST",
		"/auth/0.1/login",
		"",
		func(env controllers.Env) http.HandlerFunc { return env.Login },
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
		"LogoutToken",
		"POST",
		"/auth/0.1/logout",
		"",
		func(env controllers.Env) http.HandlerFunc { return env.LogoutToken },
	},
	Route{
		"UserData",
		"GET",
		"/data/{email}",
		"",
		func(env controllers.Env) http.HandlerFunc { return env.User },
	},
	Route{
		"UserData",
		"GET",
		"/data/{email}/{device}",
		"",
		func(env controllers.Env) http.HandlerFunc { return env.Device },
	},
	Route{
		"UserData",
		"GET",
		"/data/{email}/{device}/{feature}",
		"",
		func(env controllers.Env) http.HandlerFunc { return env.Feature },
	},
}
