package main

import (
	"context"
	"flag"
	"github.com/duongbm/greenlight-gin/internal/data"
	"github.com/duongbm/greenlight-gin/internal/jsonlog"
	"github.com/duongbm/greenlight-gin/internal/mailer"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}

// define an application struct to hold dependencies for HTTP handler, helper, middlewares, ...
type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
	mailer mailer.Mailer
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

	//SMTP config
	flag.StringVar(&cfg.smtp.host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP server host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 2525, "SMTP server port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "6618a4d7f9af60", "SMTP server username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "70ae820f1324af", "SMTP server password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Greenlight <no-reply@greenlight.duongbm.net>", "SMTP sender")

	// Initialize a new logger which write messages to the standard out stream
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	// call openDB() to create then connection pool
	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err, nil)
	}

	logger.Info("database connection pool established.", nil)

	// Declare an instance of application struct, containing the config struct and logger
	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	err = app.serve()
	// Start HTTP Server
	if err != nil {
		logger.Fatal(err, nil)
	}
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
