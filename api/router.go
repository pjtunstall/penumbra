package api

import (
	"log"
	"net/http"
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
        case http.MethodPost:
            h.SubmitLogin(w, r)
        default:
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    })

    mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
        log.Println("register")
        if r.Method == http.MethodGet {
            h.RenderRegister(w, r)
        } else if r.Method == http.MethodPost {
            log.Println("register post")
            h.SubmitRegister(w, r)
        } else {
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    })

    mux.HandleFunc(("/dashboard"), func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            h.RenderDashboard(w, r)
        } else {
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    })

    mux.HandleFunc("/tasks/create", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            h.RenderCreateTask(w, r)
        } else {
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    })

    // mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
    //     if r.Method == http.MethodPost {
    //         h.SubmitLogout(w, r)
    //     } else {
    //         http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    //     }
    // })

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
