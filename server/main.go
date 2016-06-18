package main

import (
	"net/http"
	"os"

	"github.com/gorilla/context"
	"github.com/michaelg9/ISOC/server/routers"

	"github.com/urfave/negroni"
)

func main() {
	router := routers.NewRouter()

	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.UseHandler(context.ClearHandler(router))

	port := os.Getenv("PORT")
	http.ListenAndServe(":"+port, n)
}
