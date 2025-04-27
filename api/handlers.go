package api

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"dts/app"
	"dts/db"
)

type Handler interface {
    RenderLogin(http.ResponseWriter, *http.Request)
    SubmitLogin(http.ResponseWriter, *http.Request)
    RenderRegister(http.ResponseWriter, *http.Request)
    SubmitRegister(http.ResponseWriter, *http.Request)
    HandleCreate(http.ResponseWriter, *http.Request)
    // SubmitCreateTask(http.ResponseWriter, *http.Request)
    GetTask(http.ResponseWriter, *http.Request, string)
    HandleDashboard(http.ResponseWriter, *http.Request)
    HandleHome(http.ResponseWriter, *http.Request)
}

type RealHandler struct {
    store db.Store
    templates *template.Template
}

func NewHandler(store db.Store, templates *template.Template) *RealHandler {
    return &RealHandler{store: store, templates: templates}
}

func (h *RealHandler) HandleHome(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("session_token")
    if err != nil {
        log.Println("Error getting cookie: ", err)
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    _, err = h.store.GetUserIdFromSessionToken(cookie.Value)
    if err != nil {
        log.Println("Error getting user id: ", err)
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (h *RealHandler) RenderPage(w http.ResponseWriter, r *http.Request, page string, data any) {
    pageAndOtherData := struct {
        Page string
        Data any
    }{
        Page: page,
        Data: data,
    }

    err := h.templates.ExecuteTemplate(w, "layout", pageAndOtherData)
    if err != nil {
        http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
    }
}

func (h *RealHandler) RenderLogin(w http.ResponseWriter, r *http.Request) {
    h.RenderPage(w, r, "login", nil)
}

func (h *RealHandler) RenderRegister(w http.ResponseWriter, r *http.Request) {
    h.RenderPage(w, r, "register", nil)
}

func (h *RealHandler) SubmitLogin(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        http.Error(w, "Bad Request", http.StatusBadRequest)
        log.Println("Error parsing form: ", err)
        return
    }

    user, err := h.store.GetUserByEmail(r.FormValue("email"))
    if err != nil {
        http.Error(w, "Bad Request", http.StatusBadRequest)
        log.Println("Error getting user: ", err)
        return
    }

    if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(r.FormValue("password"))); err != nil {
        http.Error(w, "Bad Request", http.StatusBadRequest)
        log.Println("Error comparing passwords: ", err)
        return
    }

    // Create session and store a hash of it in the database, in the users table. Set cookie. Fetch task titles, ids, and due dates. Redirect to `/dashboard`, which will display the list.
    sessionToken, expiresAt, err := h.store.AddSessionToken(user.Id)
    if err != nil {
        log.Println("Error adding session: ", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    http.SetCookie(w, &http.Cookie{
        Name:     "session_token",
        Value:    sessionToken,
        Path:     "/",
        HttpOnly: true,
        Secure:   false, // todo: Set to true (https) in production.
        Expires:  expiresAt,
    })

    http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (h *RealHandler) SubmitRegister(w http.ResponseWriter, r *http.Request) {
    log.Println("SubmitRegister")
    err := r.ParseForm()
    if err != nil {
        http.Error(w, "Bad Request", http.StatusBadRequest)
        return
    }

    // todo: prevent passwords larget than 72 bytes, bycrypt's limit
    // todo: 
    password_hash, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), 10)
    if err != nil {
        http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    user := app.User{
        Name: r.FormValue("name"),
        PasswordHash: password_hash,
        Email: r.FormValue("email"),
        Phone: r.FormValue("phone"),
    }

    if err := h.store.CreateUser(user); err != nil {
        log.Println("Error creating user: ", err)
        http.Error(w, "Internal Server Error: ", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func (h *RealHandler) HandleDashboard(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("session_token")
    if err != nil {
        log.Println("Error getting cookie: ", err)
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    user_id, err := h.store.GetUserIdFromSessionToken(cookie.Value)
    if err != nil {
        log.Println("Error getting user id: ", err)
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    data, err := h.store.GetAllTasks(user_id)
    if err != nil {
        log.Println("Error getting tasks: ", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    h.RenderPage(w, r, "dashboard", data)
}

func (h *RealHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("session_token")
    if err != nil {
        log.Println("Error getting cookie: ", err)
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    _, err = h.store.GetUserIdFromSessionToken(cookie.Value)
    if err != nil {
        log.Println("Error getting user id: ", err)
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    h.RenderPage(w, r, "create", nil)
}

// func (h *RealHandler) SubmitCreateTask(w http.ResponseWriter, r *http.Request) {
//     var input app.Task
//     if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
//         http.Error(w, "invalid input", http.StatusBadRequest)
//         return
//     }

//     if input.Due.IsZero() {
//         input.Due = time.Now().Add(24 * time.Hour)
//     }

//     cookie, err := r.Cookie("session_token")
//     if err != nil {
//         http.Error(w, "not logged in", http.StatusUnauthorized)
//         return
//     }

//     user_id, err := h.store.GetUserIdFromSessionToken(cookie.Value)
//     if err != nil {
//         http.Error(w, "not logged in", http.StatusUnauthorized)
//         return
//     }

//     input.UserId = user_id

//     id, err := h.store.SubmitCreateTask(input)
//     if err != nil {
//         http.Error(w, "failed to create task", http.StatusInternalServerError)
//         return
//     }

//     input.Id = id
//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(input)
// }

func (h *RealHandler) GetTask(w http.ResponseWriter, r *http.Request, id string) {
    i, err := strconv.Atoi(id)
    if err != nil {
        http.Error(w, "invalid Id", http.StatusBadRequest)
        return
    }

    task, err := h.store.GetTaskById(i)
    if err != nil {
        http.Error(w, "not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(task)
}
