package router

import (
	controller "p8ion/server/controller/auth"
)

func ApiRouter() {
	userRoutes := Router.Group("/api")
	userRoutes.GET("/oauth/google", controller.OAuthRequest)
}
