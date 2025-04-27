package api

// todo: This test file uses testify. Elsewhere, e.g. router_test, testify is not used. Decide on a consistent system. Less dependencies is good, but testify might make the code clearer due to its more declarative style.

import (
	"dts/app"
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

// func (m *MockSQLiteStore) RenderCreateTask() (int, error) {
//     args := m.Called()
//     return args.Int(0), args.Error(1)
// }

// func (m *MockSQLiteStore) SubmitCreateTask(task app.Task) (int, error) {
//     args := m.Called(task)
//     return args.Int(0), args.Error(1)
// }

func (m *MockSQLiteStore) GetTaskById(id int) (app.Task, error) {
    args := m.Called(id)
    return args.Get(0).(app.Task), args.Error(1)
}

// func TestSubmitCreateTask(t *testing.T) {
//     mockStore := &MockSQLiteStore{}
//     var templates *template.Template
//     handler := NewHandler(mockStore, templates)

//     fixedTime := time.Date(2025, 4, 21, 10, 0, 0, 0, time.UTC)
// 	task := app.Task{
//     	Title: "Test Task",
//     	Due:   fixedTime,
// 	}

//     mockStore.On("SubmitCreateTask", task).Return(1, nil)

//     body, _ := json.Marshal(task)
//     req := httptest.NewRequest("POST", "/tasks/create", bytes.NewReader(body))
//     w := httptest.NewRecorder()

//     handler.SubmitCreateTask(w, req)

//     assert.Equal(t, http.StatusOK, w.Code)
//     var createdTask app.Task
//     err := json.NewDecoder(w.Body).Decode(&createdTask)
//     assert.NoError(t, err)
//     assert.Equal(t, 1, createdTask.Id)
//     mockStore.AssertExpectations(t)
// }

// func TestGetTask(t *testing.T) {
//     mockStore := &MockSQLiteStore{}
//     var templates *template.Template
//     handler := NewHandler(mockStore, templates)

// 	due := time.Date(2025, 4, 21, 10, 0, 0, 0, time.UTC)
//     task := app.Task{
//         UserId:      1,
//         Title:       "Test Task",
//         Description: "Testing",
//         Status:      "open",
//         Done:        0,
//         Due:         due,
//     }

//     mockStore.On("GetTaskById", 1).Return(task, nil)

//     req := httptest.NewRequest("GET", "/tasks/1", nil)
//     w := httptest.NewRecorder()

//     handler.GetTask(w, req, "1")

//     assert.Equal(t, http.StatusOK, w.Code)
//     var returnedTask app.Task
//     err := json.NewDecoder(w.Body).Decode(&returnedTask)
//     assert.NoError(t, err)
//     assert.Equal(t, 1, returnedTask.Id)
//     assert.Equal(t, "Test Task", returnedTask.Title)
//     mockStore.AssertExpectations(t)
// }
