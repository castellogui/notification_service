package routes

import (
	handlers "notification_service/internal/api/handlers/notification"

	"github.com/gin-gonic/gin"
)

func RegisterNotificationRoutes(rg *gin.RouterGroup) {
	notification := rg.Group("/notifications")

	notification.POST("/", handlers.CreateNotification)
	notification.GET("/:id", nil)
	notification.PUT("/:id", nil)
	notification.DELETE("/:id", nil)
}
