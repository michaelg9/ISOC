package main

import (
	"net/http"

	"github.com/urfave/negroni"
)

func main() {
	router := NewRouter()

	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.UseHandler(router)

	http.ListenAndServe(":8080", n)
}
