package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var websocketUpgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

// Manager manage clients list
type Manager struct {
	clients ClientsList
	sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		clients: make(ClientsList),
	}
}

// serveWS upgrade http connection into websocket and runs read/write messages goroutines for each client
func (m *Manager) serveWS(w http.ResponseWriter, r *http.Request) {
	log.Println("new connection")
	wsConn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("serveWS:", err)
		return
	}
	client := NewClient(wsConn, m)
	m.addClient(client)
	go client.readMessages()
	go client.writeMessages()
}

// addClient adds client to the clients list
func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	m.clients[client] = struct{}{}
}

// removeClient removes client from the clients list
func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	if _, ex := m.clients[client]; ex {
		client.connection.Close()
		delete(m.clients, client)
	}

}
