package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type server struct {
	rooms map[string]*room
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func newServer() *server {
	return &server{
		rooms: make(map[string]*room),
	}
}

func (s *server)handleConnection(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	s.newClient(ws)
	defer ws.Close()
}



func (s *server) newClient(conn *websocket.Conn) {
	log.Printf("Client Connected: %s\n", conn.RemoteAddr())

	c := &client{
		conn: conn,
		nick: "anon",
	}

	c.read()

	log.Printf("%s | %s ", c.conn.RemoteAddr(), c.nick)

}