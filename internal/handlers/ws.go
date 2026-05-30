package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Jeno7u/server-app-course/internal/ws"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type wsMessage struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type wsResponse struct {
	Type     string `json:"type"`
	RoomID   string `json:"room_id,omitempty"`
	Username string `json:"username,omitempty"`
	Text     string `json:"text,omitempty"`
	Detail   string `json:"detail,omitempty"`
}

func ConnectRoom(c *gin.Context) {
	roomID := c.Param("room_id")
	username := c.Query("username")

	if strings.TrimSpace(username) == "" {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err == nil {
			conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1008, "Username required"))
			conn.Close()
		}
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	ws.Manager.Connect(roomID, username, conn)

	defer func() {
		ws.Manager.Disconnect(roomID, username)
		conn.Close()
	}()

	// Event about connection
	ws.Manager.Broadcast(roomID, wsResponse{
		Type:     "message",
		RoomID:   roomID,
		Username: username,
		Text:     "Joined the room",
	})

	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var payload wsMessage
		if err := json.Unmarshal(msgBytes, &payload); err == nil && payload.Type == "message" {
			if len(payload.Text) > 300 {
				conn.WriteJSON(wsResponse{
					Type:   "error",
					Detail: "Message is too long",
				})
			} else {
				ws.Manager.Broadcast(roomID, wsResponse{
					Type:     "message",
					RoomID:   roomID,
					Username: username,
					Text:     payload.Text,
				})
			}
		}
	}
}

func GetRoomUsers(c *gin.Context) {
	roomID := c.Param("room_id")
	users := ws.Manager.GetUsers(roomID)

	c.JSON(http.StatusOK, gin.H{
		"room_id": roomID,
		"users":   users,
	})
}
