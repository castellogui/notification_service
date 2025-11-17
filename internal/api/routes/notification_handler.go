package routes

import "github.com/gin-gonic/gin"

func RegisterNotificationRoutes(rg *gin.RouterGroup) {
	notification := rg.Group("/notifications")

	notification.POST("/", nil)
}
