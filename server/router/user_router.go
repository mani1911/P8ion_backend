package router

import (
	controller "p8ion/server/controller/user"
	"p8ion/server/middleware"
)

func UserRouter() {
	userRoutes := Router.Group("/user")
	userRoutes.POST("/signup", controller.SignupUser)
	userRoutes.Use(middleware.Auth)
	{
		//Dummy Protected Route
		userRoutes.GET("/user", controller.Dummy)
		userRoutes.GET("/images/:userId", controller.GetImageData) //Jeeshnu
		userRoutes.POST("/image", controller.ParseImage)           //Vezha
	}
	// userRoutes.GET("/getUser", controller.GetUserFromJwt)

}
