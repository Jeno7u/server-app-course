package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Jeno7u/server-app-course/internal/app"
	"github.com/Jeno7u/server-app-course/internal/models"
	"github.com/Jeno7u/server-app-course/internal/storage"
)

func TestDependenciesAndRouting(t *testing.T) {
	router := app.SetupRouter()

	t.Run("GET /users/me", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/users/me", nil)
		req.Header.Set("X-User-Id", "15")
		req.Header.Set("X-User-Role", "user")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %v", w.Code)
		}
	})

	t.Run("User without X-User-Id gets 401", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/users/me", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %v", w.Code)
		}
	})

	t.Run("Normal user accesses /admin/stats -> 403", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/admin/stats", nil)
		req.Header.Set("X-User-Id", "15")
		req.Header.Set("X-User-Role", "user")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Fatalf("expected 403, got %v", w.Code)
		}
	})

	t.Run("Admin accesses /admin/stats -> 200", func(t *testing.T) {
		storage.GlobalTaskStorage.Reset()
		storage.GlobalTaskStorage.Create(models.Task{Status: "todo"})
		storage.GlobalTaskStorage.Create(models.Task{Status: "done"})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/admin/stats", nil)
		req.Header.Set("X-User-Id", "1")
		req.Header.Set("X-User-Role", "admin")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %v", w.Code)
		}

		var stats struct {
			TotalTasks int            `json:"total_tasks"`
			ByStatus   map[string]int `json:"by_status"`
		}
		json.Unmarshal(w.Body.Bytes(), &stats)

		if stats.TotalTasks != 2 {
			t.Fatalf("Expected 2 tasks, got %v", stats.TotalTasks)
		}
	})

	t.Run("Admin deletes foreign task -> 204", func(t *testing.T) {
		storage.GlobalTaskStorage.Reset()
		storage.GlobalTaskStorage.Create(models.Task{OwnerID: 20})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/admin/tasks/1", nil)
		req.Header.Set("X-User-Id", "1")
		req.Header.Set("X-User-Role", "admin")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Fatalf("expected 204, got %v", w.Code)
		}
	})
}
