package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/Heiji57/ETB-BE/internal/domain"
	repository "github.com/Heiji57/ETB-BE/internal/repository/sqlc"
)

const (
	wsReadDeadline  = 60 * time.Second
	wsPingInterval  = 50 * time.Second
	wsWriteDeadline = 10 * time.Second
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Hub struct {
	mu          sync.RWMutex
	connections map[string]*websocket.Conn
}

func NewHub() *Hub {
	return &Hub{connections: make(map[string]*websocket.Conn)}
}

func (h *Hub) register(deviceID string, conn *websocket.Conn) {
	h.mu.Lock()
	h.connections[deviceID] = conn
	h.mu.Unlock()
}

func (h *Hub) unregister(deviceID string) {
	h.mu.Lock()
	delete(h.connections, deviceID)
	h.mu.Unlock()
}

type WsHandler struct {
	hub  *Hub
	repo *repository.Queries
}

func NewWs(hub *Hub, repo *repository.Queries) *WsHandler {
	return &WsHandler{hub: hub, repo: repo}
}

func (h *WsHandler) Handle(c *gin.Context) {
	deviceID := c.Query("deviceId")
	if deviceID == "" {
		domain.Fail(c, domain.ErrInvalidRequest)
		return
	}

	userID := c.GetString("userId")
	if _, err := h.repo.GetDevice(c.Request.Context(), deviceID, userID); err != nil {
		domain.Fail(c, domain.ErrNotFound)
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("ws upgrade error: %v", err)
		return
	}
	h.hub.register(deviceID, conn)
	defer func() {
		h.hub.unregister(deviceID)
		conn.Close()
	}()

	conn.SetReadDeadline(time.Now().Add(wsReadDeadline))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(wsReadDeadline))
		return nil
	})

	ticker := time.NewTicker(wsPingInterval)
	defer ticker.Stop()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, raw, err := conn.ReadMessage()
			if err != nil {
				break
			}
			var msg domain.WSMessage
			if err := json.Unmarshal(raw, &msg); err != nil {
				writeJSON(conn, map[string]string{"type": "error", "message": "invalid message format"})
				continue
			}
			switch msg.Type {
			case "location", "device_status":
				writeJSON(conn, map[string]string{"type": "ack", "refType": msg.Type})
			default:
				writeJSON(conn, map[string]string{"type": "error", "message": "unknown message type"})
			}
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(wsWriteDeadline))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func writeJSON(conn *websocket.Conn, v interface{}) {
	conn.SetWriteDeadline(time.Now().Add(wsWriteDeadline))
	data, _ := json.Marshal(v)
	conn.WriteMessage(websocket.TextMessage, data)
}
