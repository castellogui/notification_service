package notification

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) DeleteNotification(c *gin.Context) {
	id := c.Param("id")
	
	err := h.dbWriter.DeleteNotification(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete notification", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification deleted successfully"})
}