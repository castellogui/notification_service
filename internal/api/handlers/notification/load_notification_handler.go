package notification

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


func (h *Handler) GetNotification(c *gin.Context) {
	userID := c.Param("user_id")
	id := c.Param("id")
	
	notification, err := h.dbWriter.GetNotification(c.Request.Context(), id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get notification", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notification": notification})
}