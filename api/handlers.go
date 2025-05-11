package api

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"penumbra/app"
	"penumbra/db"
)

type PageAndOtherData struct {
    Page    string
    Data    any
}

type TaskView struct {
    Id        uuid.UUID
    Title     string
    Status    string
    DuePretty string
    Description string
}

func (t TaskView) String() string {
    return fmt.Sprintf("Task{Id: %s, Title: %s, Description: %s, Status: %s, Due: %s}",
        t.Id, t.Title, t.Description, t.Status, t.DuePretty)
}

type Handler interface {
    RenderLogin(http.ResponseWriter, *http.Request)
    SubmitLogin(http.ResponseWriter, *http.Request)
    RenderRegister(http.ResponseWriter, *http.Request)
    SubmitRegister(http.ResponseWriter, *http.Request)
    RenderCreateTask(http.ResponseWriter, *http.Request)
    SubmitCreateTask(http.ResponseWriter, *http.Request, int) // The `int` is the user's id.
    GetTask(http.ResponseWriter, *http.Request, uuid.UUID)
    MarkTaskDone(http.ResponseWriter, *http.Request, uuid.UUID)
    HandleDashboard(http.ResponseWriter, *http.Request)
    HandleHome(http.ResponseWriter, *http.Request)
    HandleLogout(http.ResponseWriter, *http.Request)
    HandleAllTasks(http.ResponseWriter, *http.Request, int) // The `int` is the user's id.
    DeleteTask(http.ResponseWriter, *http.Request, uuid.UUID)
    UpdateTask(http.ResponseWriter, *http.Request, uuid.UUID)
    HandleAbout(http.ResponseWriter, *http.Request)
    HandleProtected(http.ResponseWriter, *http.Request, func(http.ResponseWriter, *http.Request))
    HandleProtectedWithTaskId(http.ResponseWriter, *http.Request, func(http.ResponseWriter, *http.Request, uuid.UUID), string)
    HandleProtectedWithUserId(http.ResponseWriter, *http.Request, func(http.ResponseWriter, *http.Request, int))
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

    sessionToken, err := uuid.Parse(cookie.Value)
    if err != nil {
        http.Error(w, "not logged in", http.StatusUnauthorized)
        return
    }

    _, err = h.store.GetUserIdFromSessionToken(sessionToken)
    if err != nil {
        log.Println("Error getting user id: ", err)
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (h *RealHandler) RenderPage(w http.ResponseWriter, r *http.Request, page string, data any) {
    pageAndOtherData := PageAndOtherData{
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
        log.Println("Error parsing form: ", err)
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    user, err := h.store.GetUserByEmail(r.FormValue("email"))
    if err != nil {
        log.Println("Error getting user: ", err)
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(r.FormValue("password"))); err != nil {
        log.Println("Error comparing passwords: ", err)
        http.Redirect(w, r, "/login", http.StatusSeeOther)
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
        Value:    sessionToken.String(),
        Path:     "/",
        HttpOnly: true,
        Secure:   false, // TODO: Set to true (https) in production.
        Expires:  expiresAt,
    })

    http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (h *RealHandler) SubmitRegister(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        http.Error(w, "Bad Request", http.StatusBadRequest)
        return
    }

    password := r.FormValue("password")
    if len(password) > 72 {
        http.Error(w, "Password too long", http.StatusBadRequest)
        return
    }

    password_hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
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

    if cookie.Value == "" {
        log.Println("No session token in cookie")
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    sessionToken, err := uuid.Parse(cookie.Value)
    if err != nil {
        log.Println("Error parsing session token: ", err)
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    user_id, err := h.store.GetUserIdFromSessionToken(sessionToken)
    if err != nil {
        log.Println("Error getting user id from hashed session token: ", err)
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    preData, err := h.store.GetAllTasks(user_id)
    if err != nil {
        log.Println("Error getting tasks: ", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    sort.Slice(preData, func(i, j int) bool {
        return preData[i].Due.Before(preData[j].Due)
    })

    var data []TaskView
    for _, task := range preData {
        data = append(data, TaskView{
            Id:          task.Id,
            Title:       task.Title,
            Status:      task.Status,
            DuePretty:   task.Due.Format("Mon Jan 2 2006"),
        })
    }
   
    h.RenderPage(w, r, "dashboard", data)
}

func (h *RealHandler) RenderCreateTask(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("session_token")
    if err != nil {
        log.Println("Error getting cookie: ", err)
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    if cookie.Value == "" {
        log.Println("No session token in cookie")
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    sessionToken, err := uuid.Parse(cookie.Value)
    if err != nil {
        log.Println("Error parsing session token: ", err)
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    _, err = h.store.GetUserIdFromSessionToken(sessionToken)
    if err != nil {
        log.Println("Error getting user id: ", err)
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    h.RenderPage(w, r, "create", nil)
}

func (h *RealHandler) SubmitCreateTask(w http.ResponseWriter, r *http.Request, userId int) {
    r.ParseForm()
    title := r.FormValue("title")
	description := r.FormValue("description")
	dueStr := r.FormValue("due")

    dueDate, err := time.Parse("Mon Jan 2 2006", dueStr)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

    dueDate = time.Date(dueDate.Year(), dueDate.Month(), dueDate.Day(), 23, 59, 59, 999999999, dueDate.Location())

    task := app.Task{
        Id:          uuid.New(),
        UserId:      userId,
		Title:       title,
		Description: description,
		Due:         dueDate,
	}

    err = h.store.CreateTask(task)
    if err != nil {
        http.Error(w, "failed to create task", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (h *RealHandler) GetTask(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
    task, err := h.store.GetTaskById(id)
    if err != nil {
        http.Error(w, "not found", http.StatusNotFound)
        return
    }

    prettyTask := TaskView{
        Id:          task.Id,
        Title:       task.Title,
        Status:      task.Status,
        DuePretty:   task.Due.Format("Mon Jan 2 2006"),
        Description: task.Description,
    }

    h.RenderPage(w, r, "task", prettyTask)
}

func (h *RealHandler) HandleAbout(w http.ResponseWriter, r *http.Request) {
    h.RenderPage(w, r, "about", nil)
}

func (h *RealHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
    http.SetCookie(w, &http.Cookie{
        Name:    "session_token",
        Value:   "",
        Expires: time.Unix(0, 0),
    })
    http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *RealHandler) HandleAllTasks(w http.ResponseWriter, r *http.Request, userId int) {
    preData, err := h.store.GetAllTasks(userId)
    if err != nil {
        log.Println("Error getting tasks: ", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    sort.Slice(preData, func(i, j int) bool {
        return preData[i].Due.Before(preData[j].Due)
    })

    var data []TaskView
    for _, task := range preData {
        data = append(data, TaskView{
            Id:          task.Id,
            Title:       task.Title,
            Status:      task.Status,
            Description: task.Description,
            DuePretty:   task.Due.Format("Mon Jan 2 2006"),
        })
    }

    h.RenderPage(w, r, "tasks", data)
}

func (h *RealHandler) DeleteTask(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
    err := h.store.DeleteTask(id)
    if err != nil {
        http.Error(w, "not found", http.StatusNotFound)
        return
    }

    http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (h *RealHandler) MarkTaskDone(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
    err := h.store.SetTaskDone(id)
    if err != nil {
        http.Error(w, "not found", http.StatusNotFound)
        return
    }

    http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (h *RealHandler) UpdateTask(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
    title := r.FormValue("title")
    status := r.FormValue("status")
    description := r.FormValue("description")
    due := r.FormValue("due")

    dueDate, err := time.Parse("Mon Jan 2 2006", due)
    if err != nil {
        http.Error(w, "Invalid date format", http.StatusBadRequest)
        return
    }

    dueDate = time.Date(dueDate.Year(), dueDate.Month(), dueDate.Day(), 23, 59, 59, 999999999, dueDate.Location())

    updatedTask := app.Task{
        Id:          id,
        Title:       title,
        Status:      status,
        Description: description,
        Due:         dueDate,
    }

    err = h.store.UpdateTask(updatedTask)
    if err != nil {
        http.Error(w, "Error updating task", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (h *RealHandler) HandleProtected(w http.ResponseWriter, r *http.Request, handler func(http.ResponseWriter, *http.Request)) {
    cookie, err := r.Cookie("session_token")
    if err != nil {
        log.Println("Error getting cookie: ", err)
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    sessionToken, err := uuid.Parse(cookie.Value)
    if err != nil {
        http.Error(w, "not logged in", http.StatusUnauthorized)
        return
    }

    _, err = h.store.GetUserIdFromSessionToken(sessionToken)
    if err != nil {
        log.Println("Error getting user id: ", err)
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    handler(w, r)
}

func (h *RealHandler) HandleProtectedWithTaskId(w http.ResponseWriter, r *http.Request, handler func(http.ResponseWriter, *http.Request, uuid.UUID), idString string) {
    cookie, err := r.Cookie("session_token")
    if err != nil {
        log.Println("Error getting cookie: ", err)
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    sessionToken, err := uuid.Parse(cookie.Value)
    if err != nil {
        http.Error(w, "not logged in", http.StatusUnauthorized)
        return
    }

    _, err = h.store.GetUserIdFromSessionToken(sessionToken)
    if err != nil {
        log.Println("Error getting user id: ", err)
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    taskId, err := uuid.Parse(idString)
    if err != nil {
        http.Error(w, "invalid id", http.StatusNotFound)
        return
    }

    handler(w, r, taskId)
}

func (h *RealHandler) HandleProtectedWithUserId(w http.ResponseWriter, r *http.Request, handler func(http.ResponseWriter, *http.Request, int)) {
    cookie, err := r.Cookie("session_token")
    if err != nil {
        log.Println("Error getting cookie: ", err)
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    sessionToken, err := uuid.Parse(cookie.Value)
    if err != nil {
        http.Error(w, "not logged in", http.StatusUnauthorized)
        return
    }

    userId, err := h.store.GetUserIdFromSessionToken(sessionToken)
    if err != nil {
        log.Println("Error getting user id: ", err)
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    handler(w, r, userId)
}