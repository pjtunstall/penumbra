package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"dts/api"
	"dts/db"
)

//go:embed templates/*
var tmplFS embed.FS
var templates = template.Must(template.ParseFS(tmplFS, "templates/*.gohtml"))

func main() {
    store := db.NewSQLiteStore("data/dev.db")
    handler := api.NewHandler(store, templates)
    router := api.NewRouter(handler)

    log.Println("Server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}
