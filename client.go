package main

import (
	"github.com/gorilla/websocket"
	"log"
)

type ClientsList map[*Client]struct{}

type Client struct {
	connection *websocket.Conn
	manager    *Manager
	limit      chan []byte // limit is used to avoid concurrent writes to WS connection
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		limit:      make(chan []byte),
	}
}

func (c *Client) readMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()
	for {
		msgType, msg, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println(err)
			}
			break
		}
		log.Println(msgType)
		log.Println(string(msg))
	}
}

func (c *Client) writeMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()
	for {
		select {
		case msg, ok := <-c.limit:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed", err)
				}
				return
			}
			if err := c.connection.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Println("failed to send message:", err)
			}
			log.Println("message sent")
		}

	}
}
