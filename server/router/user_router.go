package router

import (
	controller "p8ion/server/controller/user"
)

func UserRouter() {
	userRoutes := Router.Group("/user")
	userRoutes.POST("/signup", controller.SignupUser)
	userRoutes.GET("/user", controller.Dummy)
}
