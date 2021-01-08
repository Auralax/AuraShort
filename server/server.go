package server

import (
	"aurashort/server/data"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

// StartServer is used to start the http server and initialize  the routing
func StartServer() {
	r := mux.NewRouter()

	r.HandleFunc("/create", data.CreateLink).Methods("POST") // Used to create a new random short url

	if os.Getenv("CUSTOM_LINKS_ENABLED") == "true" {
		r.HandleFunc("/custom", data.CreateCustomLink).Methods("POST") // Used to create a custom short url
	}
	r.HandleFunc("/{id}", data.Redirect).Methods("GET") // Used to redirect to a url by id

	err := http.ListenAndServe(":8000", r)
	if err != nil {
		log.Fatal(err)
	}
}
