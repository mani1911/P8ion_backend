package utils

import "github.com/gin-gonic/gin"

func SendResponse(c *gin.Context, statusCode int, message interface{}) {
	c.JSON(statusCode, gin.H{
		"status_code": statusCode,
		"message":     message,
	})
	c.Abort()
}
