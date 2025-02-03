package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type client struct {
	socket   *websocket.Conn
	send     chan *SendMessage
	room     *room
	userData map[string]string
}

func (c *client) read() {
	defer c.socket.Close()

	msgs, err := GetMessages()
	if err != nil {
		log.Println("Failed to get messages: ", err)
	} else {
		for _, msg := range msgs {
			c.send <- &msg
		}
	}

	for {
		var msg SendMessage
		if err := c.socket.ReadJSON(&msg); err != nil {
			break
		}

		msg.CreatedAt = time.Now()
		msg.UserName = c.userData["name"]
		msg.Email = c.userData["email"]
		msg.AvatarURL = c.userData["avatar_url"]
		c.room.forward <- &msg

		SaveMessage(Message{
			UserID:    c.userData["user_id"],
			Message:   msg.Message,
			CreatedAt: msg.CreatedAt,
		})
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
