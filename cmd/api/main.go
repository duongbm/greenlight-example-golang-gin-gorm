package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/duongbm/greenlight-gin/internal/data"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	db   struct {
		dsn         string
		maxConn     int
		MaxIdleConn int
		maxIdleTime string
	}
}

// define an application struct to hold dependencies for HTTP handler, helper, middlewares, ...
type application struct {
	config config
	logger *log.Logger
	models data.Models
}

func main() {
	var cfg config

	// Read value of port and env command-line flags into config struct
	flag.IntVar(&cfg.port, "port", 8000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "environment(development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("DB_DSN"), "PostgreSQL connection DSN")
	flag.IntVar(&cfg.db.maxConn, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.MaxIdleConn, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max idle timeout")
	flag.Parse()

	// Initialize a new logger which write messages to the standard out stream
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// call openDB() to create then connection pool
	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Printf("database connection pool established.")

	// Declare an instance of application struct, containing the config struct and logger
	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf("localhost:%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Minute * 10,
		WriteTimeout: time.Minute * 30,
	}

	// Start HTTP Server
	logger.Printf("Starting %s server on %d", cfg.env, cfg.port)
	err = srv.ListenAndServe()
	logger.Fatal(err)
}

func openDB(cfg config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.db.dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// create a context with 5 second timeout deadline
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.db.maxConn)
	sqlDB.SetMaxIdleConns(cfg.db.MaxIdleConn)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	sqlDB.SetConnMaxIdleTime(duration)

	err = sqlDB.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
