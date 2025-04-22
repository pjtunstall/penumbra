package api

import (
	"encoding/json"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"dts/app"
	"dts/db"
)

type Handler interface {
    RenderHome(http.ResponseWriter, *http.Request)
    RenderRegister(http.ResponseWriter, *http.Request)
    SubmitRegister(http.ResponseWriter, *http.Request)
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
        Page string
    }{
        Flag: flag,
        Page: "home",
    }

    err := h.templates.ExecuteTemplate(w, "layout", data)
    if err != nil {
        http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
    }
}

func (h *RealHandler) RenderRegister(w http.ResponseWriter, r *http.Request) {
    data := struct {
        Page string
    }{
        Page: "register",
    }

    err := h.templates.ExecuteTemplate(w, "layout", data)
    if err != nil {
        http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
    }
}

func (h *RealHandler) SubmitRegister(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        http.Error(w, "Bad Request", http.StatusBadRequest)
        return
    }

    user := app.User{
        Name: r.FormValue("name"),
        PasswordHash: r.FormValue("password"),
        Email: r.FormValue("email"),
        Phone: r.FormValue("phone"),
    }

    if err := h.store.CreateUser(user); err != nil {
        http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/login", http.StatusSeeOther)
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
