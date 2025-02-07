package websocket

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"

	"github.com/fhiroki/chat/internal/domain/message"
	"github.com/fhiroki/chat/internal/trace"
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

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := r.upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	userData := make(map[string]string)
	if authCookie, err := req.Cookie("auth"); err == nil {
		userData = parseUserData(authCookie.Value)
	}

	client := newClient(socket, r, userData, r.msgService)
	r.Join <- client
	defer func() { r.Leave <- client }()
	go client.write()
	client.read()
}

func parseUserData(cookieValue string) map[string]string {
	userData := make(map[string]string)
	if decoded, err := base64.StdEncoding.DecodeString(cookieValue); err == nil {
		if err := json.Unmarshal(decoded, &userData); err != nil {
			log.Printf("Failed to unmarshal auth cookie: %v", err)
		}
	}
	return userData
}
