package websocketmanager

import (
	"sync"

	"github.com/gorilla/websocket"
)

type ConnectionManager struct {
	mu sync.RWMutex
	activeConnections map[string][]*websocket.Conn
}

var manager *ConnectionManager
var once sync.Once

func GetManager() *ConnectionManager {
	once.Do(func() {
		manager = &ConnectionManager{
			activeConnections: make(map[string][]*websocket.Conn),
		}
	})
	return manager
}

func (cm *ConnectionManager) Connect(conn *websocket.Conn, userID string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.activeConnections[userID] = append(cm.activeConnections[userID], conn)
}

func (cm *ConnectionManager) Disconnect(conn *websocket.Conn, userID string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	conns := cm.activeConnections[userID]
	for i, c := range conns {
		if c == conn {
			_ = c.Close()

			cm.activeConnections[userID] = append(conns[:i], conns[i+1:]...)
			break
		}
	}

	if len(cm.activeConnections[userID]) == 0 {
		delete(cm.activeConnections, userID)
	}
}

func (cm *ConnectionManager) SendPersonalMessage(message string, userID string) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	for _, conn := range cm.activeConnections[userID] {
		_ = conn.WriteMessage(websocket.TextMessage, []byte(message))
	}
}

func (cm *ConnectionManager) Broadcast(message string) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	for _, conns := range cm.activeConnections {
		for _, conn := range conns {
			_ = conn.WriteMessage(websocket.TextMessage, []byte(message))
		}
	}
}