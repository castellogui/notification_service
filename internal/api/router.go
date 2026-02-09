package api

import (
	handlers "notification_service/internal/api/handlers/notification"
	"notification_service/internal/api/routes"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"github.com/gocql/gocql"
	"notification_service/internal/infra"
)

func SetupRouter(r *gin.Engine, w *kafka.Writer, s *gocql.Session) {
	v1 := r.Group("/v1")

	scyllaWriter := infra.NewScyllaWriter(s)
	notificationHandler := handlers.NewHandler(w, scyllaWriter)
	routes.RegisterNotificationRoutes(v1, notificationHandler)
}