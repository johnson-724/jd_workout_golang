package auth

import (
	"github.com/gin-gonic/gin"
)

func InfoAction(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ok",
		"data": gin.H{
			"id": uid,
		},
	})
}
