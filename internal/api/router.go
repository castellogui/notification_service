package api

import (
	"notification_service/internal/api/routes"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	v1 := r.Group("/v1")

	routes.RegisterNotificationRoutes(v1)
}
