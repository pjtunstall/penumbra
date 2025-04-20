package main

import (
	"log"
	"net/http"

	"dts/backend/internal/api"
	"dts/backend/internal/db"
)

func main() {
    store := db.NewSQLiteStore("data/dev.db")
    handler := api.NewHandler(store)
    router := api.NewRouter(handler)

    log.Println("Server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}
