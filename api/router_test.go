package api_test

// import (
// 	"net/http"
// 	"testing"
// )

// type mockHandler struct {
// 	homeCalled bool
// 	renderRegisterCalled bool
// 	submitRegisterCalled bool
// 	submitLoginCalled bool
// 	renderDashboardCalled bool
// 	getTaskCalled    bool
// 	gotId        string
// 	RenderCreateTaskCalled bool
// 	SubmitCreateTaskCalled bool
// }

// // todo: Test this too.
// func (m *mockHandler) RenderHome(w http.ResponseWriter, r *http.Request) {
// 	m.homeCalled = true
// 	w.WriteHeader(http.StatusOK)
// }

// func (m *mockHandler) RenderRegister(w http.ResponseWriter, r *http.Request) {
// 	m.renderRegisterCalled = true
// 	w.WriteHeader(http.StatusOK)
// }

// func (m *mockHandler)SubmitRegister(w http.ResponseWriter, r *http.Request) {
// 	m.submitRegisterCalled = true
// 	w.WriteHeader(http.StatusOK)
// }

// func (m *mockHandler) SubmitLogin(w http.ResponseWriter, r *http.Request) {
// 	m.submitLoginCalled = true
// 	w.WriteHeader(http.StatusOK)
// }

// func (m *mockHandler) RenderDashboard(w http.ResponseWriter, r *http.Request) {
// 	m.renderDashboardCalled = true
// 	w.WriteHeader(http.StatusOK)
// }

// func (m *mockHandler) RenderCreateTask(w http.ResponseWriter, r *http.Request) {
// 	m.RenderCreateTaskCalled = true
// 	w.WriteHeader(http.StatusOK)
// }

// // func (m *mockHandler) SubmitCreateTask(w http.ResponseWriter, r *http.Request) {
// // 	m.submitCreateTaskCalled = true
// // 	w.WriteHeader(http.StatusCreated)
// // }

// func (m *mockHandler) GetTask(w http.ResponseWriter, r *http.Request, id string) {
// 	m.getTaskCalled = true
// 	m.gotId = id
// 	w.WriteHeader(http.StatusOK)
// }

// func TestNewRouter(t *testing.T) {
// 	// mock := &mockHandler{}
// 	// router := api.NewRouter(mock)

// 	// t.Run("GET /tasks/create calls RenderCreateTask", func(t *testing.T) {
// 	// 	req := httptest.NewRequest(http.MethodPost, "/tasks", nil)
// 	// 	rec := httptest.NewRecorder()

// 	// 	router.ServeHTTP(rec, req)

// 	// 	if !mock.submitCreateTaskCalled {
// 	// 		t.Errorf("RenderCreateTask was not called")
// 	// 	}
// 	// 	if rec.Code != http.StatusCreated {
// 	// 		t.Errorf("expected status 201, got %d", rec.Code)
// 	// 	}
// 	// })

// 	// t.Run("POST /tasks/create calls SubmitCreateTask", func(t *testing.T) {
// 	// 	req := httptest.NewRequest(http.MethodPost, "/tasks", nil)
// 	// 	rec := httptest.NewRecorder()

// 	// 	router.ServeHTTP(rec, req)

// 	// 	if !mock.SubmitCreateTaskCalled {
// 	// 		t.Errorf("SubmitCreateTask was not called")
// 	// 	}
// 	// 	if rec.Code != http.StatusCreated {
// 	// 		t.Errorf("expected status 201, got %d", rec.Code)
// 	// 	}
// 	// })

// 	// t.Run("GET /tasks/123 calls GetTask with correct Id", func(t *testing.T) {
// 	// 	req := httptest.NewRequest(http.MethodGet, "/tasks/123", nil)
// 	// 	rec := httptest.NewRecorder()

// 	// 	router.ServeHTTP(rec, req)

// 	// 	if !mock.getTaskCalled {
// 	// 		t.Errorf("GetTask was not called")
// 	// 	}
// 	// 	if mock.gotId != "123" {
// 	// 		t.Errorf("expected Id '123', got '%s'", mock.gotId)
// 	// 	}
// 	// 	if rec.Code != http.StatusOK {
// 	// 		t.Errorf("expected status 200, got %d", rec.Code)
// 	// 	}
// 	// })

// 	// t.Run("GET /tasks returns 405", func(t *testing.T) {
// 	// 	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
// 	// 	rec := httptest.NewRecorder()

// 	// 	router.ServeHTTP(rec, req)

// 	// 	if rec.Code != http.StatusMethodNotAllowed {
// 	// 		t.Errorf("expected status 405, got %d", rec.Code)
// 	// 	}
// 	// })

// 	// t.Run("POST /tasks/123 returns 405", func(t *testing.T) {
// 	// 	req := httptest.NewRequest(http.MethodPost, "/tasks/123", nil)
// 	// 	rec := httptest.NewRecorder()

// 	// 	router.ServeHTTP(rec, req)

// 	// 	if rec.Code != http.StatusMethodNotAllowed {
// 	// 		t.Errorf("expected status 405, got %d", rec.Code)
// 	// 	}
// 	// })
// }
