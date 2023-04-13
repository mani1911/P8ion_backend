package router

import (
	"p8ion/server/middleware"

	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

func Init() {
	Router = gin.Default()
	Router.Use(middleware.LoggerMiddleware)
	ApiRouter()
	UserRouter()
}
