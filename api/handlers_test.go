package api

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"penumbra/app"
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

func (m *MockSQLiteStore) AddSessionToken(user_id int) (uuid.UUID, time.Time, error) {
    args := m.Called(user_id)
    return args.Get(0).(uuid.UUID), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockSQLiteStore) GetUserIdFromSessionToken(sessionToken uuid.UUID) (int, error) {
    args := m.Called(sessionToken)
    return args.Int(0), args.Error(1)
}

func (m *MockSQLiteStore) GetAllTasks(user_id int) ([]app.Task, error) {
    args := m.Called(user_id)
    return args.Get(0).([]app.Task), args.Error(1)
}

func (m *MockSQLiteStore) RenderCreateTask() (int, error) {
    args := m.Called()
    return args.Int(0), args.Error(1)
}

func (m *MockSQLiteStore) SubmitCreateTask(task app.Task) (int, error) {
    args := m.Called(task)
    return args.Int(0), args.Error(1)
}

func (m *MockSQLiteStore) GetTaskById(id uuid.UUID) (app.Task, error) {
    args := m.Called(id)
    return args.Get(0).(app.Task), args.Error(1)
}

func (m *MockSQLiteStore) MarkTaskDone(id int) error {
    args := m.Called(id)
    return args.Error(0)
}

func (m *MockSQLiteStore) CreateTask(task app.Task) error {
    args := m.Called(task)
    return args.Error(0)
}

func (m *MockSQLiteStore) DeleteTask(id uuid.UUID) error {
    args := m.Called(id)
    return args.Error(0)
}

func (m *MockSQLiteStore) UpdateTask(task app.Task) error {
    args := m.Called(task)
    return args.Error(0)
}

func (m *MockSQLiteStore) SetTaskDone(id uuid.UUID) error {
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

func TestHandleHomeValidCookie(t *testing.T) {
	mockStore := &MockSQLiteStore{}
	handler := &RealHandler{store: mockStore}

	req, err := http.NewRequest("GET", "/home", nil)
	if err != nil {
		t.Fatal(err)
	}

	cookie := &http.Cookie{Name: "session_token", Value: uuid.New().String()}
	req.AddCookie(cookie)

    mockStore.On("GetUserIdFromSessionToken", mock.AnythingOfType("uuid.UUID")).Return(1, nil).Once()

	rr := httptest.NewRecorder()

	handler.HandleHome(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Fatalf("Expected status %d but got %d", http.StatusSeeOther, rr.Code)
	}

	if location := rr.Header().Get("Location"); location != "/dashboard" {
		t.Fatalf("Expected redirect to /dashboard but got %s", location)
	}
}

func TestSubmitLogin(t *testing.T) {
	password := "password123"
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to generate password hash: %v", err)
	}
	
	t.Logf("Generated password hash: %s", string(passwordHash))
	
	err = bcrypt.CompareHashAndPassword(passwordHash, []byte(password))
	if err != nil {
		t.Fatalf("Generated hash fails verification: %v", err)
	}
	
	mockStore := new(MockSQLiteStore)
	handler := &RealHandler{
		store: mockStore,
	}
	
	mockUser := app.User{
		Id:           1,
		Email:        "test@example.com",
		PasswordHash: passwordHash,
	}
	
	mockStore.On("GetUserByEmail", "test@example.com").Return(mockUser, nil).Once()
	
	mockSessionToken := uuid.New()
	mockStore.On("AddSessionToken", mockUser.Id).Return(mockSessionToken, time.Now().Add(time.Hour), nil).Once()
	
	form := url.Values{}
	form.Add("email", "test@example.com")
	form.Add("password", password)
	
	req := httptest.NewRequest("POST", "/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	rr := httptest.NewRecorder()
	handler.SubmitLogin(rr, req)
	
	t.Logf("Response status code: %d", rr.Code)
	t.Logf("Response headers: %v", rr.Header())
	
	assert.Equal(t, http.StatusSeeOther, rr.Code)
	
	redirectLocation := rr.Header().Get("Location")
	t.Logf("Redirected to: %s", redirectLocation)
	assert.Equal(t, "/dashboard", redirectLocation)
	
	cookies := rr.Result().Cookies()
	assert.Len(t, cookies, 1, "Expected exactly one cookie")
	
	if len(cookies) > 0 {
		cookie := cookies[0]
		t.Logf("Cookie: Name=%s, Value=%s, Path=%s", cookie.Name, cookie.Value, cookie.Path)
		assert.Equal(t, "session_token", cookie.Name)
		assert.Equal(t, mockSessionToken.String(), cookie.Value)
	}
	
	mockStore.AssertExpectations(t)
}