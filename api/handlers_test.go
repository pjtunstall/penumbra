package api

// todo: This test file uses testify. Elsewhere, e.g. router_test, testify is not used. Decide on a consistent system. Less dependencies is good, but testify might make the code clearer due to its more declarative style.

import (
	"bytes"
	"dts/app"
	"encoding/json"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSQLiteStore struct {
    mock.Mock
}

func (m *MockSQLiteStore) CreateUser(user app.User) error {
    args := m.Called(user)
    return args.Error(0)
}

func (m *MockSQLiteStore) CreateTask(task app.Task) (int, error) {
    args := m.Called(task)
    return args.Int(0), args.Error(1)
}

func (m *MockSQLiteStore) GetTask(id int) (app.Task, error) {
    args := m.Called(id)
    return args.Get(0).(app.Task), args.Error(1)
}

func TestCreateTask(t *testing.T) {
    mockStore := &MockSQLiteStore{}
    var templates *template.Template
    handler := NewHandler(mockStore, templates)

    fixedTime := time.Date(2025, 4, 21, 10, 0, 0, 0, time.UTC)
	task := app.Task{
    	Title: "Test Task",
    	Due:   fixedTime,
	}

    mockStore.On("CreateTask", task).Return(1, nil)

    body, _ := json.Marshal(task)
    req := httptest.NewRequest("POST", "/tasks", bytes.NewReader(body))
    w := httptest.NewRecorder()

    handler.CreateTask(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    var createdTask app.Task
    err := json.NewDecoder(w.Body).Decode(&createdTask)
    assert.NoError(t, err)
    assert.Equal(t, 1, createdTask.ID)
    mockStore.AssertExpectations(t)
}

func TestGetTask(t *testing.T) {
    mockStore := &MockSQLiteStore{}
    var templates *template.Template
    handler := NewHandler(mockStore, templates)

	fixedTime := time.Date(2025, 4, 21, 10, 0, 0, 0, time.UTC)
    task := app.Task{
        ID:    1,
        Title: "Test Task",
        Due:   fixedTime,
    }

    mockStore.On("GetTask", 1).Return(task, nil)

    req := httptest.NewRequest("GET", "/tasks/1", nil)
    w := httptest.NewRecorder()

    handler.GetTask(w, req, "1")

    assert.Equal(t, http.StatusOK, w.Code)
    var returnedTask app.Task
    err := json.NewDecoder(w.Body).Decode(&returnedTask)
    assert.NoError(t, err)
    assert.Equal(t, 1, returnedTask.ID)
    assert.Equal(t, "Test Task", returnedTask.Title)
    mockStore.AssertExpectations(t)
}
