package ws
package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

type RoomManager struct {
	sync.RWMutex
	rooms map[string]map[string]*websocket.Conn
}

var Manager = &RoomManager{
	rooms: make(map[string]map[string]*websocket.Conn),
}

func (m *RoomManager) Connect(roomID, username string, conn *websocket.Conn) {
	m.Lock()
	defer m.Unlock()
	if m.rooms[roomID] == nil {
		m.rooms[roomID] = make(map[string]*websocket.Conn)
	}
	m.rooms[roomID][username] = conn
}

func (m *RoomManager) Disconnect(roomID, username string) {
	m.Lock()
	defer m.Unlock()
	if m.rooms[roomID] != nil {
		delete(m.rooms[roomID], username)
		if len(m.rooms[roomID]) == 0 {
			delete(m.rooms, roomID)
		}
	}
}

func (m *RoomManager) Broadcast(roomID string, message interface{}) {
	m.RLock()
	defer m.RUnlock()
	if room, ok := m.rooms[roomID]; ok {
		for _, conn := range room {
			// Ignore error in broadcast for simplicity
			_ = conn.WriteJSON(message)
		}
	}
}

func (m *RoomManager) GetUsers(roomID string) []string {
	m.RLock()
	defer m.RUnlock()
	users := make([]string, 0)
	if room := m.rooms[roomID]; room != nil {
		for uname := range room {
			users = append(users, uname)
		}
	}
	return users
}
