package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hellojonas/barracuda/internal/barracuda"
	"github.com/hellojonas/barracuda/pkg/logs"
	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
)

func main() {
    r := chi.NewRouter()
    r.Use(middleware.Logger)

    db, err := newDB()

    if err != nil {
	log.Fatal(err)
    }

    b := barracuda.New(db)
    logDir := os.Getenv("HOME") + "/.barracuda/logs"
    logger, err := newLogger(logDir)

    if err != nil {
	log.Fatalf("error creating logger %v", err)
    }

    defer logger.Close()

    b.SetLogger(logger)
    bRouter := barracuda.Router(b)

    c := cron.New()
    _, err = c.AddFunc("0 4-23/6 * * *", func() {
	logger.Info("job to refresh articles has been started")
	b.RefreshArticles()
    })

    if err != nil {
	fmt.Println(err)
	logger.Error("failed to register job to refresh articles. %v", err)
	os.Exit(1)
    }

    c.Start()

    r.Mount("/api", bRouter)

    logger.Info("Application started, listenning on port :8080")
    err = http.ListenAndServe(":8080", r)

    if err != nil {
	logger.Error("application exit with error")
    }
}

func newDB() (*sql.DB, error) {
    datasource := os.Getenv("DATABASE_URL")

    if datasource == "" {
	return nil, errors.New("datasource not set")
    }

    db, err := sql.Open("postgres", datasource)

    if err != nil {
	return nil, fmt.Errorf("error opening database connection. %v", err)
    }

    if err = db.Ping(); err != nil {
	return nil, fmt.Errorf("could not establish connection to database %v", err)
    }
    return db, nil
}

func newLogger(dir string) (*logs.BLogger, error) {
    stat, err := os.Stat(dir)

    if (err != nil && os.IsNotExist(err)) || !stat.IsDir() {
	err = nil
	err = os.MkdirAll(dir, fs.ModePerm)
    }

    if err != nil {
	return nil, err
    }

    logFile := dir + "/barracuda.log"
    file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_APPEND, fs.ModePerm)

    if os.IsNotExist(err) {
	err = nil
	file, err = os.Create(logFile)
    }

    if err != nil {
	return nil, fmt.Errorf("error opening log file. %v", err)
    }

    logger := logs.New(file)
    return logger, nil
}
