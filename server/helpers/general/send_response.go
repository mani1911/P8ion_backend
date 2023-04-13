package helper

import (
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func SendResponse(c *gin.Context, status int, res interface{}) {
	c.JSON(status, res)
	c.Abort()
}

func SendError(c *gin.Context, status int, message string) {
	SendResponse(c, status, ErrorResponse{Message: message})
}
