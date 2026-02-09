package notification

import (
	"encoding/json"
	"net/http"

	"notification_service/internal/pusher/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"time"
)

type Handler struct {
	writer *kafka.Writer
}

func NewHandler(w *kafka.Writer) *Handler {
	return &Handler{writer: w}
}

func (h *Handler) CreateNotification(c *gin.Context) {
	var envelope domain.Envelope
	if err := c.ShouldBindJSON(&envelope); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	envelope.CreatedAt = &now

	msg, err := json.Marshal(envelope)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal envelope"})
		return
	}

	if envelope.DeliverAt != nil && envelope.DeliverAt.After(time.Now()) {
		c.JSON(http.StatusAccepted, gin.H{"status": "notification_scheduled", "notification": envelope})
		return
	}

	if err := h.writer.WriteMessages(c.Request.Context(), kafka.Message{
		Key:   []byte(uuid.New().String()),
		Value: msg,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to produce message"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "notification_queued", "notification": envelope})
}
