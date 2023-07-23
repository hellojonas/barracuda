package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hellojonas/barracuda/internal/barracuda"
)


func main() {
    r := chi.NewRouter()
    r.Use(middleware.Logger)
    r.Mount("/api", barracuda.Routes())

    fmt.Println("Application listenning on port :8080")
    http.ListenAndServe(":8080", r)
}
