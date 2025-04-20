package api

import (
	"embed"
	"encoding/json"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"dts/backend/internal/app"
	"dts/backend/internal/db"
)

//go:embed templates/*
var tmplFS embed.FS
var templates = template.Must(template.ParseFS(tmplFS, "templates/*.html"))

type TaskHandler interface {
    RenderHome(http.ResponseWriter, *http.Request)
    CreateTask(http.ResponseWriter, *http.Request)
    GetTask(http.ResponseWriter, *http.Request, string)
}

type Handler struct {
    store db.Store
}

func NewHandler(store db.Store) *Handler {
    return &Handler{store: store}
}

func (h *Handler) RenderHome(w http.ResponseWriter, r *http.Request) {
    err := templates.ExecuteTemplate(w, "index.html", nil)
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}


func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
    var input app.Task
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "invalid input", http.StatusBadRequest)
        return
    }

    if input.Due.IsZero() {
        input.Due = time.Now().Add(24 * time.Hour)
    }

    id, err := h.store.CreateTask(input)
    if err != nil {
        http.Error(w, "failed to create task", http.StatusInternalServerError)
        return
    }

    input.ID = id
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(input)
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request, id string) {
    i, err := strconv.Atoi(id)
    if err != nil {
        http.Error(w, "invalid ID", http.StatusBadRequest)
        return
    }

    task, err := h.store.GetTask(i)
    if err != nil {
        http.Error(w, "not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(task)
}
