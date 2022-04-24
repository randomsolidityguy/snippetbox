package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
}


func main() {
	// Command-line flag
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Parse command-line flag
	flag.Parse()

	// Custom loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Initialize new instance of application struct containing dependencies
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
	}

	// Starting server
	srv := &http.Server {
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
	}

	infoLog.Printf("Starting server on port %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}