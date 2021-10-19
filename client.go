package main

import (
	"github.com/gorilla/websocket"
	"log"
	"strings"
)

type client struct {
	conn *websocket.Conn
	nick string
	rooms map[string]*room
	commands chan<- command
}

func (c *client) read() {
	for {
		_, message, err := c.conn.ReadMessage()

		if err != nil {
			log.Println(err)
			return
		}

		var msg  = string(message)

		msg = strings.Replace(msg,"\n", "", -1)
		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

		switch cmd {
		case "/nick":
			c.commands <- command{
				id: CMD_NICK,
				client: c,
				args: args,
			}
		case "/join":
			c.commands <- command{
				id: CMD_JOIN,
				client: c,
				args: args,
			}
		case "/rooms":
			c.commands <- command{
				id: CMD_ROOMS,
				client: c,
				args: args,
			}
		case "/msg":
			c.commands <- command{
				id: CMD_MSG,
				client: c,
				args: args,
			}
		case "/leave":
			c.commands <- command{
				id: CMD_LEAVE,
				client: c,
				args: args,
			}
		default:
			c.sendMessage("Error, unknown command")
		}
	}
}

func (c *client)sendMessage(message string)  {
	err := c.conn.WriteMessage(websocket.TextMessage, []byte(message))

	if err != nil {
		log.Print(err)
	}
}