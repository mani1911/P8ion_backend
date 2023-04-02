package middleware

import (
	"time"

	"p8ion/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggerMiddleware(c *gin.Context) {
	begin := time.Now()

	c.Next()

	end := time.Now()

	status := c.Writer.Status()

	var level logrus.Level = logrus.InfoLevel

	if status >= 400 {
		level = logrus.WarnLevel
	}

	if status >= 500 {
		level = logrus.ErrorLevel
	}

	utils.Logger.WithFields(logrus.Fields{
		"ip":       c.ClientIP(),
		"method":   c.Request.Method,
		"path":     c.Request.URL.Path,
		"status":   status,
		"duration": end.Sub(begin),
	}).Log(level, "Request")
}
