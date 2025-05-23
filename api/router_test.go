package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"penumbra/api"
)

type MockHandler struct {
	mock.Mock
}

func (m *MockHandler) HandleHome(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *MockHandler) RenderLogin(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *MockHandler) SubmitLogin(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *MockHandler) RenderRegister(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *MockHandler) SubmitRegister(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *MockHandler) HandleProtected(w http.ResponseWriter, r *http.Request, handlerFunc func(http.ResponseWriter, *http.Request)) {
	m.Called(w, r, handlerFunc)
	handlerFunc(w, r)
}

func (m *MockHandler) HandleProtectedWithTaskId(w http.ResponseWriter, r *http.Request, handlerFunc func(http.ResponseWriter, *http.Request, uuid.UUID), id string) {
	m.Called(w, r, handlerFunc, id)
	handlerFunc(w, r, uuid.MustParse(id))
}

func (m *MockHandler) HandleProtectedWithUserId(w http.ResponseWriter, r *http.Request, handlerFunc func(http.ResponseWriter, *http.Request, int)) {
	m.Called(w, r, handlerFunc)
	handlerFunc(w, r, 1)
}

func (m *MockHandler) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *MockHandler) HandleAbout(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *MockHandler) RenderCreateTask(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *MockHandler) SubmitCreateTask(w http.ResponseWriter, r *http.Request, userId int) {
	m.Called(w, r, userId)
}

func (m *MockHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *MockHandler) HandleAllTasks(w http.ResponseWriter, r *http.Request, userId int) {
	m.Called(w, r, userId)
}

func (m *MockHandler) GetTask(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	m.Called(w, r, id)
}

func (m *MockHandler) MarkTaskDone(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	m.Called(w, r, id)
}

func (m *MockHandler) UpdateTask(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	m.Called(w, r, id)
}

func (m *MockHandler) DeleteTask(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	m.Called(w, r, id)
}

func TestDashboardRoute(t *testing.T) {
	mockHandler := new(MockHandler)
	router := api.NewRouter(mockHandler)

	mockHandler.On("HandleProtected", mock.Anything, mock.Anything, mock.Anything).Maybe().Run(func(args mock.Arguments) {
		fn := args.Get(2).(func(http.ResponseWriter, *http.Request))
		fn(args.Get(0).(http.ResponseWriter), args.Get(1).(*http.Request))
	})
	mockHandler.On("HandleDashboard", mock.Anything, mock.Anything).Maybe()

	mockHandler.On("HandleProtectedWithTaskId", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().
	Run(func(args mock.Arguments) {
		fn := args.Get(2).(func(http.ResponseWriter, *http.Request, uuid.UUID))
		id := uuid.MustParse(args.Get(3).(string))
		fn(args.Get(0).(http.ResponseWriter), args.Get(1).(*http.Request), id)
	})

	mockHandler.On("HandleProtectedWithUserId", mock.Anything, mock.Anything, mock.Anything).Maybe().
	Run(func(args mock.Arguments) {
		fn := args.Get(2).(func(http.ResponseWriter, *http.Request, int))
		fn(args.Get(0).(http.ResponseWriter), args.Get(1).(*http.Request), 1)
	})

	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockHandler.AssertExpectations(t)
}

func TestNewRouter_Routes(t *testing.T) {
	mockHandler := new(MockHandler)
	router := api.NewRouter(mockHandler)

	type testCase struct {
		name       string
		method     string
		url        string
		body       []byte
		expectFunc func()
		expectCode int
	}

	cases := []testCase{
		{
			name:   "Home GET",
			method: http.MethodGet,
			url:    "/",
			expectFunc: func() {
				mockHandler.On("HandleHome", mock.Anything, mock.Anything).Once()
			},
			expectCode: http.StatusOK,
		},
		{
			name:   "Login POST",
			method: http.MethodPost,
			url:    "/login",
			expectFunc: func() {
				mockHandler.On("SubmitLogin", mock.Anything, mock.Anything).Once()
			},
			expectCode: http.StatusOK,
		},
		{
			name:   "Register GET",
			method: http.MethodGet,
			url:    "/register",
			expectFunc: func() {
				mockHandler.On("RenderRegister", mock.Anything, mock.Anything).Once()
			},
			expectCode: http.StatusOK,
		},
		{
			name:   "Task Done POST",
			method: http.MethodPost,
			url:    "/tasks/done/123e4567-e89b-12d3-a456-426614174000",
			body:   mustJSON(map[string]bool{"checked": true}),
			expectFunc: func() {
				id := "123e4567-e89b-12d3-a456-426614174000"
				mockHandler.On(
					"HandleProtectedWithTaskId",
					mock.Anything,
					mock.Anything,
					mock.AnythingOfType("func(http.ResponseWriter, *http.Request, uuid.UUID)"),
					id,
				).Once()
			
				mockHandler.On("MarkTaskDone", mock.Anything, mock.Anything, uuid.MustParse(id)).Once()
			},
			expectCode: http.StatusOK,
		},			
		{
			name:   "Logout GET",
			method: http.MethodGet,
			url:    "/logout",
			expectFunc: func() {
				mockHandler.On("HandleLogout", mock.Anything, mock.Anything).Once()
			},
			expectCode: http.StatusOK,
		},
		
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			var req *http.Request
			if tc.body != nil {
				req = httptest.NewRequest(tc.method, tc.url, bytes.NewReader(tc.body))
			} else {
				req = httptest.NewRequest(tc.method, tc.url, nil)
			}

			tc.expectFunc()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectCode, rec.Code)
			mockHandler.AssertExpectations(t)
		})
	}
}

func mustJSON(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}
