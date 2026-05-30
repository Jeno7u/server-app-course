package handlers
package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/Jeno7u/server-app-course/internal/handlers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gin-gonic/gin"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	memUsers := router.Group("/mem_users")
	{
		memUsers.POST("", handlers.CreateMemUser)
		memUsers.GET("/:id", handlers.GetMemUser)
		memUsers.DELETE("/:id", handlers.DeleteMemUser)
	}

	return router
}

func TestCreateUser(t *testing.T) {
	handlers.ResetMemDB()
	router := setupTestRouter()

	userIn := handlers.MemUserIn{
		Username: gofakeit.Username(),
		Age:      gofakeit.Number(18, 99),
	}
	body, _ := json.Marshal(userIn)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/mem_users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected status code 201, got %v", w.Code)
	}

	var out handlers.MemUserOut
	json.Unmarshal(w.Body.Bytes(), &out)

	if out.Username != userIn.Username || out.Age != userIn.Age {
		t.Fatalf("Returned data does not match struct: expected %v, got %v", userIn, out)
	}
}

func TestGetUser(t *testing.T) {
	handlers.ResetMemDB()
	router := setupTestRouter()

	// Seed User
	userIn := handlers.MemUserIn{
		Username: gofakeit.Username(),
		Age:      gofakeit.Number(18, 99),
	}
	body, _ := json.Marshal(userIn)
	wResp := httptest.NewRecorder()
	reqReq, _ := http.NewRequest("POST", "/mem_users", bytes.NewBuffer(body))
	reqReq.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(wResp, reqReq)

	var createdUser handlers.MemUserOut
	json.Unmarshal(wResp.Body.Bytes(), &createdUser)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/mem_users/"+strconv.Itoa(createdUser.ID), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code 200, got %v", w.Code)
	}

	var out handlers.MemUserOut
	json.Unmarshal(w.Body.Bytes(), &out)
	if out.ID != createdUser.ID {
		t.Fatalf("Returned incorrect user ID")
	}
}

func TestGetMissingUser(t *testing.T) {
	handlers.ResetMemDB()
	router := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/mem_users/999", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("Expected status 404, got %v", w.Code)
	}
}

func TestDeleteUser(t *testing.T) {
	handlers.ResetMemDB()
	router := setupTestRouter()

	// Seed User
	userIn := handlers.MemUserIn{
		Username: gofakeit.Username(),
		Age:      gofakeit.Number(18, 99),
	}
	body, _ := json.Marshal(userIn)
	wResp := httptest.NewRecorder()
	reqReq, _ := http.NewRequest("POST", "/mem_users", bytes.NewBuffer(body))
	reqReq.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(wResp, reqReq)

	var createdUser handlers.MemUserOut
	json.Unmarshal(wResp.Body.Bytes(), &createdUser)

	// Delete
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/mem_users/"+strconv.Itoa(createdUser.ID), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("Expected status 204, got %v", w.Code)
	}

	// Verify not found
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/mem_users/"+strconv.Itoa(createdUser.ID), nil)
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusNotFound {
		t.Fatalf("Expected status 404 after deletion, got %v", w2.Code)
	}
}

func TestDeleteMissingUser(t *testing.T) {
	handlers.ResetMemDB()
	router := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/mem_users/999", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("Expected status 404, got %v", w.Code)
	}
}
