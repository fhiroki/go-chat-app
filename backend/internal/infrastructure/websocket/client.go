package websocket

import (
	"context"
	"log"
	"time"

	"github.com/fhiroki/chat/internal/domain/message"
	"github.com/gorilla/websocket"
)

type client struct {
	socket         *websocket.Conn
	send           chan *message.Message
	room           *room
	userData       map[string]string
	messageService message.MessageService
}

func newClient(socket *websocket.Conn, room *room, userData map[string]string, messageService message.MessageService) *client {
	return &client{
		socket:         socket,
		send:           make(chan *message.Message),
		room:           room,
		userData:       userData,
		messageService: messageService,
	}
}

func (c *client) read() {
	defer c.socket.Close()

	ctx := context.Background()
	messages, err := c.messageService.FindAll(ctx)
	if err != nil {
		log.Println("Failed to get messages: ", err)
	} else {
		for _, msg := range messages {
			msg.AttachUserData(c.userData)
			c.send <- msg
		}
	}

	for {
		var msg message.Message
		if err := c.socket.ReadJSON(&msg); err != nil {
			break
		}

		msg.UserID = c.userData["user_id"]
		// TODO: 設定しなくても現在時刻を入れるようにしたい
		msg.CreatedAt = time.Now()
		msg.UpdatedAt = msg.CreatedAt

		if err := c.messageService.Create(ctx, &msg); err != nil {
			log.Printf("Failed to save message: %v", err)
			continue
		}

		c.room.Forward <- &msg
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
