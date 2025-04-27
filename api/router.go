package api

import (
	"net/http"
)

func NewRouter(h Handler) http.Handler {
    mux := http.NewServeMux()

    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            h.HandleHome(w, r)
        } else {
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    })

    mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            h.RenderLogin(w, r)
        case http.MethodPost:
            h.SubmitLogin(w, r)
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
        if r.Method == http.MethodGet {
            h.HandleProtected(w, r, h.HandleDashboard)
        } else {
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    })

    mux.HandleFunc("/tasks/create", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            h.HandleProtected(w, r, h.HandleCreate)
        } else {
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    })

    mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            h.HandleLogout(w, r)
        } else {
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    })

    // mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
    //     if r.Method == http.MethodPost {
    //         h.CreateTask(w, r)
    //     } else {
    //         http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    //     }
    // })

    // mux.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
    //     id := strings.TrimPrefix(r.URL.Path, "/tasks/")
    //     if r.Method == http.MethodGet {
    //         h.GetTask(w, r, id)
    //     } else {
    //         http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    //     }
    // })

    return mux
}
