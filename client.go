package main

import (
	"time"

	"github.com/gorilla/websocket"
)

type client struct {
	socket   *websocket.Conn
	send     chan *message
	room     *room
	userData map[string]string
}

func (c *client) read() {
	defer c.socket.Close()

	for {
		var msg message
		if err := c.socket.ReadJSON(&msg); err != nil {
			break
		}
		msg.When = time.Now()
		msg.Name = c.userData["name"]
		msg.Email = c.userData["email"]
		msg.AvatarURL, _ = c.room.avatar.GetAvatarURL(c)
		c.room.forward <- &msg
	}
}

func (c *client) write() {
	defer c.socket.Close()

	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
}
