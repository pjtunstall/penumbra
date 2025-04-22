package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"dts/api"
)

type mockHandler struct {
	homeCalled bool
	renderRegisterCalled bool
	submitRegisterCalled bool
	createCalled bool
	getCalled    bool
	gotID        string
}

// todo: Test this too.
func (m *mockHandler) RenderHome(w http.ResponseWriter, r *http.Request) {
	m.homeCalled = true
	w.WriteHeader(http.StatusOK)
}

func (m *mockHandler) RenderRegister(w http.ResponseWriter, r *http.Request) {
	m.renderRegisterCalled = true
	w.WriteHeader(http.StatusOK)
}

func (m *mockHandler)SubmitRegister(w http.ResponseWriter, r *http.Request) {
	m.submitRegisterCalled = true
	w.WriteHeader(http.StatusOK)
}

func (m *mockHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	m.createCalled = true
	w.WriteHeader(http.StatusCreated)
}

func (m *mockHandler) GetTask(w http.ResponseWriter, r *http.Request, id string) {
	m.getCalled = true
	m.gotID = id
	w.WriteHeader(http.StatusOK)
}

func TestNewRouter(t *testing.T) {
	mock := &mockHandler{}
	router := api.NewRouter(mock)

	t.Run("POST /tasks calls CreateTask", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/tasks", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if !mock.createCalled {
			t.Errorf("CreateTask was not called")
		}
		if rec.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d", rec.Code)
		}
	})

	t.Run("GET /tasks/123 calls GetTask with correct ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/tasks/123", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if !mock.getCalled {
			t.Errorf("GetTask was not called")
		}
		if mock.gotID != "123" {
			t.Errorf("expected ID '123', got '%s'", mock.gotID)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rec.Code)
		}
	})

	t.Run("GET /tasks returns 405", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status 405, got %d", rec.Code)
		}
	})

	t.Run("POST /tasks/123 returns 405", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/tasks/123", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status 405, got %d", rec.Code)
		}
	})
}
