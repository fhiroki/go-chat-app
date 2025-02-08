package handler

import (
	"strconv"

	"github.com/fhiroki/chat/internal/domain/message"
	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	messageService message.MessageService
}

func NewMessageHandler(messageService message.MessageService) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
	}
}

func (h *MessageHandler) GetMessages(c *gin.Context) {
	messages, err := h.messageService.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, messages)
}

func (h *MessageHandler) CreateMessage(c *gin.Context) {
	var msg message.Message
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.messageService.Create(c.Request.Context(), &msg); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Status(201)
}

func (h *MessageHandler) DeleteMessage(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid message ID"})
		return
	}

	if err := h.messageService.Delete(ctx, id); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Status(204)
}
