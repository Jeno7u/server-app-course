package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Jeno7u/server-app-course/internal/app"
	"github.com/gorilla/websocket"
)

func TestWebSocketRooms(t *testing.T) {
	router := app.SetupRouter()
	server := httptest.NewServer(router)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	t.Run("Connect to room and message", func(t *testing.T) {
		url := wsURL + "/ws/rooms/python?username=alice"
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			t.Fatalf("Failed to connect: %v", err)
		}
		defer conn.Close()

		// Read join message
		var msg map[string]interface{}
		conn.ReadJSON(&msg)
		if msg["text"] != "Joined the room" {
			t.Fatalf("Expected join msg, got %v", msg)
		}

		err = conn.WriteJSON(map[string]string{"type": "message", "text": "Всем привет"})
		if err != nil {
			t.Fatalf("WriteJSON failed: %v", err)
		}

		conn.ReadJSON(&msg)
		if msg["text"] != "Всем привет" {
			t.Fatalf("Expected text 'Всем привет', got %v", msg)
		}

		// Check HTTP route
		resp, _ := http.Get(server.URL + "/rooms/python/users")
		var roomData struct {
			Users []string `json:"users"`
		}
		json.NewDecoder(resp.Body).Decode(&roomData)

		found := false
		for _, u := range roomData.Users {
			if u == "alice" {
				found = true
			}
		}
		if !found {
			t.Fatalf("Expected alice in users, got %v", roomData.Users)
		}
	})

	t.Run("Error on long message", func(t *testing.T) {
		url := wsURL + "/ws/rooms/python2?username=bob"
		conn, _, _ := websocket.DefaultDialer.Dial(url, nil)
		defer conn.Close()

		// Skip join message
		var msg map[string]interface{}
		conn.ReadJSON(&msg)

		longText := strings.Repeat("A", 301)
		conn.WriteJSON(map[string]string{"type": "message", "text": longText})

		conn.ReadJSON(&msg)
		if msg["type"] != "error" {
			t.Fatalf("Expected error for long msg, got %v", msg)
		}
	})

	t.Run("Disconnect removes user", func(t *testing.T) {
		url := wsURL + "/ws/rooms/go?username=test"
		conn, _, _ := websocket.DefaultDialer.Dial(url, nil)
		conn.Close()

		time.Sleep(100 * time.Millisecond) // Let backend process close

		resp, _ := http.Get(server.URL + "/rooms/go/users")
		var roomData struct {
			Users []string `json:"users"`
		}
		json.NewDecoder(resp.Body).Decode(&roomData)

		if len(roomData.Users) > 0 {
			t.Fatalf("Expected empty room after disconnect, got %v", roomData.Users)
		}
	})
}
