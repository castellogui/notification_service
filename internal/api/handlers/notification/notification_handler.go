package notification

import (
	"notification_service/internal/pusher/interfaces"
	"github.com/segmentio/kafka-go"
)

type Handler struct {
	writer *kafka.Writer
	dbWriter interfaces.Writer
}

func NewHandler(w *kafka.Writer, dbWriter interfaces.Writer) *Handler {
	return &Handler{writer: w, dbWriter: dbWriter}
}