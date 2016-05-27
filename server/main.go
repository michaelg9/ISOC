package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", Index)
	log.Fatal(http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", nil))
}
