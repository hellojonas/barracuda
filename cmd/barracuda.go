package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hellojonas/barracuda/internal/barracuda"
	_ "github.com/lib/pq"
)

func main() {
    r := chi.NewRouter()
    r.Use(middleware.Logger)


    datasource := os.Getenv("DATABASE_URL")

    if datasource == "" {
	log.Fatal("datasource not set")
    }

    db, err := sql.Open("postgres", datasource)

    if err != nil {
	log.Fatalf("error opening database connection. %v", err)
    }

    if err = db.Ping(); err != nil {
	log.Fatalf("could not establish connection to database %v", err)
    }

    b := barracuda.New(db)
    bRouter := barracuda.Router(b)

    r.Mount("/api", bRouter)

    fmt.Println("Application listenning on port :8080")
    http.ListenAndServe(":8080", r)
}
