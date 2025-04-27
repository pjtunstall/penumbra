package api

import (
	"net/http"
	"strings"
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
            h.HandleProtected(w, r, h.RenderCreate)
        } else if r.Method == http.MethodPost {
            h.HandleProtected(w, r, h.SubmitCreate)
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

    mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            h.HandleAllTasks(w, r)
        } else {
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    })

    mux.HandleFunc("/task/", func(w http.ResponseWriter, r *http.Request) {
        id := strings.TrimPrefix(r.URL.Path, "/task/")
        if r.Method == http.MethodGet {
            h.GetTask(w, r, id)
        } else {
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    })
    
    // mux.HandleFunc("/task/edit/", func(w http.ResponseWriter, r *http.Request) {
    //     id := strings.TrimPrefix(r.URL.Path, "/task/edit/")
    //     if r.Method == http.MethodPost {
    //         h.EditTask(w, r, id)
    //     } else {
    //         http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    //     }
    // })

    mux.HandleFunc("/tasks/delete/", func(w http.ResponseWriter, r *http.Request) {
        id := strings.TrimPrefix(r.URL.Path, "/tasks/delete/")
        if r.Method == http.MethodPost {
            h.DeleteTask(w, r, id)
        } else {
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    })

    return mux
}
