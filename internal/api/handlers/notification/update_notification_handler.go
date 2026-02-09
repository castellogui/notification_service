package notification

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var allowedFields = map[string]bool{
	"title":     true,
	"body":      true,
	"category":  true,
	"deep_link": true,
	"data":      true,
	"read":      true,
}

func (h *Handler) UpdateNotification(c *gin.Context) {
	userID := c.Param("user_id")
	id := c.Param("id")

	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	fields := make(map[string]interface{})
	for key, val := range body {
		if allowedFields[key] {
			fields[key] = val
		}
	}

	if len(fields) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no valid fields to update"})
		return
	}

	notification, err := h.dbWriter.UpdateNotification(c.Request.Context(), userID, id, fields)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update notification", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notification": notification})
}
