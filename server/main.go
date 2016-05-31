package main

import (
	"net/http"

	"github.com/michaelg9/ISOC/server/routers"

	"github.com/urfave/negroni"
)

func main() {
	router := routers.NewRouter()

	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.UseHandler(router)

	http.ListenAndServe(":8080", n)
}
