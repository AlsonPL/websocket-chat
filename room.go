package main

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type room struct {
	name string
	clients map[*websocket.Conn]*client
}

func (r *room)broadcastSystemMessage(s *client, message string)  {
	for conn, c := range r.clients {
		if conn != s.conn {
			c.sendMessage(fmt.Sprintf("[%s] - %s", r.name, message))
		}
	}
}

func (r *room)broadcastUserMessage(s *client, message string)  {
	for _, c := range r.clients {
		c.sendMessage(fmt.Sprintf("[%s] - %s: %s", r.name, s.nick, message))
	}
}