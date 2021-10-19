package main

import "github.com/gorilla/websocket"

type room struct {
	name string
	clients map[*websocket.Conn]*client
}
