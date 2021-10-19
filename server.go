package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
)

type server struct {
	rooms map[string]*room
	commands chan command
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func newServer() *server {
	return &server{
		rooms: make(map[string]*room),
		commands: make(chan command),
	}
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NICK:
			s.nick(cmd.client, cmd.args[1])
		case CMD_JOIN:
			s.join(cmd.client, cmd.args[1])
		case CMD_ROOMS:
			s.listRooms(cmd.client)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_LEAVE:
			s.leave(cmd.client, cmd.args[1])
		}
	}
}

func (s *server)handleConnection(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	s.newClient(ws)
}

func (s *server) newClient(conn *websocket.Conn) {
	log.Printf("Client Connected: %s\n", conn.RemoteAddr())

	c := &client{
		conn: conn,
		nick: "anon",
		rooms: make(map[string]*room),
		commands: s.commands,
	}

	c.read()
}

func (s *server) nick(c *client, newNick string) {
	c.nick = newNick
	c.sendMessage("You changed your nick to: " + c.nick)
}

func (s *server) join(c *client, roomName string) {
	r, ok := s.rooms[roomName]

	if !ok {
		r = &room{
			name: roomName,
			clients: make(map[*websocket.Conn]*client),
		}
		s.rooms[roomName] = r
	}

	if _, ok := c.rooms[roomName]; ok {
		c.sendMessage("You already joined this channel")
		return
	}

	r.clients[c.conn] = c
	c.rooms[roomName] = r

	c.sendMessage(fmt.Sprintf("- You joined the %s channel -", r.name))
	r.broadcastSystemMessage(c, fmt.Sprintf("%s joined", c.nick))
}

func (s *server) msg(c *client, args []string) {
	roomName := args[1]
	msg := strings.Join(args[2:], " ")

	_, ok := c.rooms[roomName]

	if !ok {
		c.sendMessage("You must join the room before sending a message ")
		return
	}

	c.rooms[roomName].broadcastUserMessage(c, msg)
}

func (s *server) listRooms(c *client) {
	var rooms []string

	for name := range s.rooms {
		rooms = append(rooms, name)
	}

	c.sendMessage(fmt.Sprintf("Rooms: %s", strings.Join(rooms, ", ")))
}

func (s *server) leave(c *client, roomName string)  {
	room, ok := s.rooms[roomName]

	if !ok {
		c.sendMessage("You cannot leave this room")
		return
	}

	delete(s.rooms[room.name].clients, c.conn)
	delete(c.rooms, room.name)

	room.broadcastSystemMessage(c, fmt.Sprintf("%s - has left the room", c.nick))

	if len(room.clients) == 0 {
		delete(s.rooms, room.name)
		log.Printf("Deleting %s room", room.name)
	}
}