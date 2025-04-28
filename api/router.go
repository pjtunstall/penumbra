package api

import (
	"encoding/json"
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

    mux.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            h.HandleProtected(w, r, h.HandleAbout)
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

    mux.HandleFunc("/task/done/", func(w http.ResponseWriter, r *http.Request) {
        id := strings.TrimPrefix(r.URL.Path, "/task/done/")
    
        if r.Method == http.MethodPost {
            var requestBody struct {
                Checked bool `json:"checked"`
            }
    
            if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
                http.Error(w, "Invalid request body", http.StatusBadRequest)
                return
            }

            if !requestBody.Checked {
                return
            }
    
            h.DoneTask(w, r, id)
            return
        }
    
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
    })
    
    
    mux.HandleFunc("/tasks/update/", func(w http.ResponseWriter, r *http.Request) {
        id := strings.TrimPrefix(r.URL.Path, "/tasks/update/")
        if r.Method == http.MethodPost {
            h.UpdateTask(w, r, id)
        } else {
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    })

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
