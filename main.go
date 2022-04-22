package main

import (
	"fmt"
	"log"
	"net/http"
)

// Handler
func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Snippetbox"))
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	
	port := 4000
	log.Printf("Starting server on: %d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	log.Fatal(err)

	
}