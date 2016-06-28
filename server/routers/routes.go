package routers

import (
	"net/http"

	"github.com/michaelg9/ISOC/server/controllers"
)

// Route struct for creating routes
type Route struct {
	Name        string
	Method      string           // HTTP method
	Pattern     string           // URI for this route
	HandlerFunc http.HandlerFunc // Handler function specified in controllers
}

// Routes is an array of routes
type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		controllers.Index,
	},
	Route{
		"LoginWeb",
		"GET",
		"/login",
		controllers.LoginWeb,
	},
	Route{
		"Dashboard",
		"GET",
		"/dashboard",
		controllers.Dashboard,
	},
	Route{
		"Login",
		"POST",
		"/auth/0.1/login",
		controllers.Login,
	},
	Route{
		"Logout",
		"POST",
		"/auth/0.1/logout",
		controllers.Logout,
	},
	Route{
		"SignUp",
		"POST",
		"/signup",
		controllers.SignUp,
	},
	Route{
		"Upload",
		"POST",
		"/app/0.1/upload",
		controllers.Upload,
	},
	Route{
		"Download",
		"GET",
		"/data/0.1/q",
		controllers.Download,
	},
}
