package api

import (
	"log"
	"net/http"
	"strings"
)

func NewRouter(h Handler) http.Handler {
    mux := http.NewServeMux()

    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            h.RenderHome(w, r)
        } else {
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    })

    mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            // render login form
        case http.MethodPost:
            // authenticate, then:
            http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
        default:
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    })

    mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            h.RenderRegister(w, r)
        } else if r.Method == http.MethodPost {
            h.SubmitRegister(w, r)
        } else {
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    })

    mux.HandleFunc(("/dashboard"), func(w http.ResponseWriter, r *http.Request) {
        log.Println("Dashboard hit with method:", r.Method)
        if r.Method == http.MethodGet {
            h.RenderDashboard(w, r)
        } else {
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    })

    mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            h.CreateTask(w, r)
        } else {
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    })

    mux.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
        id := strings.TrimPrefix(r.URL.Path, "/tasks/")
        if r.Method == http.MethodGet {
            h.GetTask(w, r, id)
        } else {
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    })

    return mux
}
