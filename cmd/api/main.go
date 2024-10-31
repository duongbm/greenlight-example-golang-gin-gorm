package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Declare a string contain the application version number
// it should be generated automatically at build time
const version = "1.0.0"

// define config struct to hold all configuration settings for application
type config struct {
	port int
	env  string
}

// define an application struct to hold dependencies for HTTP handler, helper, middlewares, ...
type application struct {
	config config
	logger *log.Logger
}

func main() {
	var cfg config

	// Read value of port and env command-line flags into config struct
	flag.IntVar(&cfg.port, "port", 8000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "environment(development|staging|production)")
	flag.Parse()

	// Initialize a new logger which write messages to the standard out stream
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Declare an instance of application struct, containing the config struct and logger
	app := &application{
		config: cfg,
		logger: logger,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Minute * 10,
		WriteTimeout: time.Minute * 30,
	}

	// Start HTTP Server
	logger.Printf("Starting %s server on %d", cfg.env, cfg.port)
	err := srv.ListenAndServe()
	logger.Fatal(err)
}
