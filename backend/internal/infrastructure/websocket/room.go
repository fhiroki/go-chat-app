package websocket

import (
	"encoding/base64"
	"encoding/json"
	"log"

	"github.com/fhiroki/chat/internal/domain/message"
	"github.com/fhiroki/chat/internal/trace"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type room struct {
	Forward    chan *message.Message
	Join       chan *client
	Leave      chan *client
	Clients    map[*client]bool
	Tracer     trace.Tracer
	upgrader   *websocket.Upgrader
	msgService message.MessageService
}

func NewRoom(messageService message.MessageService) *room {
	return &room{
		Forward:    make(chan *message.Message),
		Join:       make(chan *client),
		Leave:      make(chan *client),
		Clients:    make(map[*client]bool),
		Tracer:     trace.Off(),
		upgrader:   &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024},
		msgService: messageService,
	}
}

func (r *room) Run() {
	for {
		select {
		case client := <-r.Join:
			r.Clients[client] = true
			r.Tracer.Trace("New client joined")
		case client := <-r.Leave:
			delete(r.Clients, client)
			close(client.send)
			r.Tracer.Trace("Client left")
		case msg := <-r.Forward:
			r.Tracer.Trace("Message received: ", msg.Content)
			for client := range r.Clients {
				client.send <- msg
				r.Tracer.Trace(" -- sent to client")
			}
		}
	}
}

func (r *room) HandleWebSocket(c *gin.Context) {
	socket, err := r.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	authCookie, err := c.Cookie("auth")
	if err != nil {
		log.Fatal("Failed to get auth cookie:", err)
		return
	}

	userData := parseUserData(authCookie)
	client := newClient(socket, r, userData, r.msgService)

	r.Join <- client
	defer func() { r.Leave <- client }()
	go client.write()
	client.read()
}

func parseUserData(cookieValue string) map[string]string {
	decoded, err := base64.StdEncoding.DecodeString(cookieValue)
	if err != nil {
		log.Fatal("Failed to decode cookie value:", err)
		return nil
	}

	var userData map[string]string
	if err = json.Unmarshal(decoded, &userData); err != nil {
		log.Fatal("Failed to decode user data JSON:", err)
		return nil
	}

	return userData
}
