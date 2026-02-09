package api

import (
	handlers "notification_service/internal/api/handlers/notification"
	"notification_service/internal/api/routes"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

func SetupRouter(r *gin.Engine, w *kafka.Writer) {
	v1 := r.Group("/v1")

	notificationHandler := handlers.NewHandler(w)
	routes.RegisterNotificationRoutes(v1, notificationHandler)
}
