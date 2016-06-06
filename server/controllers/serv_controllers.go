package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
)

// ServeCSS serves the css files static
func ServeCSS(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// TODO: Check if non-empty
	http.ServeFile(w, r, "static/"+vars["filename"])
}

// ServeIMG serves the image files in static/img
func ServeIMG(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// TODO: Check if non-empty
	http.ServeFile(w, r, "static/img/"+vars["name"])
}
