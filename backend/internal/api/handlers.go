package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"dts/backend/internal/app"
	"dts/backend/internal/db"
)

type TaskHandler interface {
    CreateTask(http.ResponseWriter, *http.Request)
    GetTask(http.ResponseWriter, *http.Request, string)
}

type Handler struct {
    store db.Store
}

func NewHandler(store db.Store) *Handler {
    return &Handler{store: store}
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
