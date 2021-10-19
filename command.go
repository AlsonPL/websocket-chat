package main

type commandID int

const (
	CMD_JOIN commandID = iota
	CMD_NICK
	CMD_ROOMS
	CMD_MSG
	CMD_LEAVE
)

type command struct {
	id commandID
	client *client
	args []string
}
