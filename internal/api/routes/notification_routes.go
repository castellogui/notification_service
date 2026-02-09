package routes

import (
	handlers "notification_service/internal/api/handlers/notification"

	"github.com/gin-gonic/gin"
)

func RegisterNotificationRoutes(rg *gin.RouterGroup, h *handlers.Handler) {
	notification := rg.Group("/notifications")

	notification.POST("/", h.CreateNotification)
	notification.GET("/:user_id/:id", h.GetNotification)
	notification.PATCH("/:user_id/:id", h.UpdateNotification)
	notification.DELETE("/:id", h.DeleteNotification)
}
