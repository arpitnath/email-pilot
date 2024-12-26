package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func EmailHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Email Handler is working"})
}
