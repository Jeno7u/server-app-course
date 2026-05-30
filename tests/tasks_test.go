package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Jeno7u/server-app-course/internal/app"
	"github.com/Jeno7u/server-app-course/internal/models"
	"github.com/Jeno7u/server-app-course/internal/storage"
)

func TestTasksAPI(t *testing.T) {
	router := app.SetupRouter()

	t.Run("Create Task Success", func(t *testing.T) {
		storage.GlobalTaskStorage.Reset()
		body := `{"title":"Подготовить тесты","description":"Написать интеграционные тесты","status":"todo","priority":4}`

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/tasks", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-Id", "10")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %v", w.Code)
		}
	})

	t.Run("Create Task Short Title 422", func(t *testing.T) {
		storage.GlobalTaskStorage.Reset()
		body := `{"title":"A","status":"todo","priority":4}`

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/tasks", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-Id", "10")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected 422, got %v", w.Code)
		}
	})

	t.Run("Unauthorized 401 Missing Header", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/tasks", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %v", w.Code)
		}
	})

	t.Run("User sees only own tasks & Filtering", func(t *testing.T) {
		storage.GlobalTaskStorage.Reset()

		// Create for User 10
		storage.GlobalTaskStorage.Create(models.Task{Title: "T1", OwnerID: 10, Status: "todo", Priority: 5})
		storage.GlobalTaskStorage.Create(models.Task{Title: "T2", OwnerID: 10, Status: "done", Priority: 1})
		// Create for User 20
		storage.GlobalTaskStorage.Create(models.Task{Title: "T3", OwnerID: 20})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/tasks?status=todo&min_priority=3", nil)
		req.Header.Set("X-User-Id", "10")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %v", w.Code)
		}

		var tasks []models.Task
		json.Unmarshal(w.Body.Bytes(), &tasks)

		if len(tasks) != 1 || tasks[0].Title != "T1" {
			t.Fatalf("Filtering failed, got: %v", tasks)
		}
	})

	t.Run("Success update status & delete", func(t *testing.T) {
		storage.GlobalTaskStorage.Reset()
		storage.GlobalTaskStorage.Create(models.Task{Title: "T1", OwnerID: 10, Status: "todo", Priority: 5})

		body := `{"status":"done"}`
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/api/tasks/1/status", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-Id", "10")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %v", w.Code)
		}

		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest("DELETE", "/api/tasks/1", nil)
		req2.Header.Set("X-User-Id", "10")
		router.ServeHTTP(w2, req2)

		if w2.Code != http.StatusNoContent {
			t.Fatalf("expected 204, got %v", w2.Code)
		}
	})

	t.Run("404 Foreign task", func(t *testing.T) {
		storage.GlobalTaskStorage.Reset()
		storage.GlobalTaskStorage.Create(models.Task{Title: "T1", OwnerID: 20, Status: "todo", Priority: 5})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/tasks/1", nil)
		req.Header.Set("X-User-Id", "10")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %v", w.Code)
		}
	})
}
