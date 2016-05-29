package main

import "net/http"

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
		Index,
	},
	Route{
		"Login",
		"GET",
		"/auth/0.1/login?uname={uname}&pwd={pwd}",
		Login,
	},
	Route{
		"Logout",
		"GET",
		"/auth/0.1/logout",
		Logout,
	},
	Route{
		"Upload",
		"POST",
		"/auth/0.1/upload",
		Upload,
	},
	Route{
		"Download",
		"GET",
		"/data/0.1/q?appid={apikey}",
		Download,
	},
}
