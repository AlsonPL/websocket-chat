package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

type client struct {
	conn *websocket.Conn
	nick string
	rooms map[*room]room
}

func (c *client) read() {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(string(message))
	}
}