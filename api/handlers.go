package api

import (
	"encoding/json"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"dts/app"
	"dts/db"
)

type Handler interface {
    RenderHome(http.ResponseWriter, *http.Request)
    CreateTask(http.ResponseWriter, *http.Request)
    GetTask(http.ResponseWriter, *http.Request, string)
}

type RealHandler struct {
    store db.Store
    templates *template.Template
}

func NewHandler(store db.Store, templates *template.Template) *RealHandler {
    return &RealHandler{store: store, templates: templates}
}

func (h *RealHandler) RenderHome(w http.ResponseWriter, r *http.Request) {
    rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
    flag := rnd.Intn(2)
    data := struct {
        Flag int
    }{
        Flag: flag,
    }

    log.Println(data.Flag)

    err := h.templates.ExecuteTemplate(w, "home.html", data)
    if err != nil {
        http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
    }
}

func (h *RealHandler) Register(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    err := h.templates.ExecuteTemplate(w, "register.html", nil)
    if err != nil {
        http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
    }
}

func (h *RealHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
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

func (h *RealHandler) GetTask(w http.ResponseWriter, r *http.Request, id string) {
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
