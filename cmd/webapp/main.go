package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"penumbra/api"
	"penumbra/db"
)

//go:embed templates/*
var tmplFS embed.FS
var templates = template.Must(template.ParseFS(tmplFS, "templates/*.html"))

func main() {
    store, err := db.NewSQLiteStore("data/dev.db")
	if err != nil {
		log.Fatalf("NewSQLiteStore failed: %v", err)
	}

    handler := api.NewHandler(store, templates)
    router := api.NewRouter(handler)

    log.Println("Server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}
