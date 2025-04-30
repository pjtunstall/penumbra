package api

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"penumbra/app"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockSQLiteStore struct {
    mock.Mock
}

func (m *MockSQLiteStore) CreateUser(user app.User) error {
    args := m.Called(user)
    return args.Error(0)
}

func (m *MockSQLiteStore) GetUserByEmail(email string) (app.User, error) {
    args := m.Called(email)
    return args.Get(0).(app.User), args.Error(1)
}

func (m *MockSQLiteStore) AddSessionToken(user_id int) (string, time.Time, error) {
    args := m.Called(user_id)
    return args.String(0), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockSQLiteStore) GetUserIdFromSessionToken(sessionToken string) (int, error) {
    args := m.Called(sessionToken)
    return args.Int(0), args.Error(1)
}

func (m *MockSQLiteStore) GetAllTasks(user_id int) ([]app.Task, error) {
    args := m.Called(user_id)
    return args.Get(0).([]app.Task), args.Error(1)
}

func (m *MockSQLiteStore) RenderCreate() (int, error) {
    args := m.Called()
    return args.Int(0), args.Error(1)
}

func (m *MockSQLiteStore) SubmitCreate(task app.Task) (int, error) {
    args := m.Called(task)
    return args.Int(0), args.Error(1)
}

func (m *MockSQLiteStore) GetTaskById(id int) (app.Task, error) {
    args := m.Called(id)
    return args.Get(0).(app.Task), args.Error(1)
}

func (m *MockSQLiteStore) DoneTask(id int) error {
    args := m.Called(id)
    return args.Error(0)
}

func (m *MockSQLiteStore) DeleteTask(id int) error {
    args := m.Called(id)
    return args.Error(0)
}

func (m *MockSQLiteStore) UpdateTask(task app.Task) error {
    args := m.Called(task)
    return args.Error(0)
}

func (m *MockSQLiteStore) SetTaskDone(id int) error {
    args := m.Called()
    return args.Error(1)
}

func TestRenderPage(t *testing.T) {
    tmpl := template.Must(template.New("layout").Parse("<html>{{.Page}}</html>"))
    handler := &RealHandler{
        templates: tmpl,
    }

    req := httptest.NewRequest(http.MethodGet, "/page", nil)
    res := httptest.NewRecorder()

    handler.RenderPage(res, req, "home", nil)

    body := res.Body.String()
    if !strings.Contains(body, "home") {
        t.Errorf("Expected page content to include 'home', but got %s", body)
    }
}


func TestHandleHomeMissingCookie(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/home", nil)
    res := httptest.NewRecorder()

    handler := &RealHandler{
    }

    handler.HandleHome(res, req)

    if res.Code != http.StatusSeeOther {
        t.Errorf("Expected status %v, got %v", http.StatusSeeOther, res.Code)
    }

    location := res.Header().Get("Location")
    if location != "/login" {
        t.Errorf("Expected redirect to /login, got %v", location)
    }
}
