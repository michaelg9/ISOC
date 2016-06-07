package routers

import (
	"net/http"

	"github.com/michaelg9/ISOC/server/controllers"
)

// Route struct for creating routes
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
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
		"ServeCSS",
		"GET",
		"/static/{filename}",
		controllers.ServeCSS,
	},
	Route{
		"ServeIMG",
		"GET",
		"/static/img/{name}",
		controllers.ServeIMG,
	},
	Route{
		"Login",
		"GET",
		"/auth/0.1/login",
		controllers.Login,
	},
	Route{
		"Logout",
		"GET",
		"/auth/0.1/logout",
		controllers.Logout,
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
