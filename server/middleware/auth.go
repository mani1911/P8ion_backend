package middleware

import (
	"fmt"
	"net/http"
	userHelper "p8ion/server/helpers/auth"
	generalHelper "p8ion/server/helpers/general"

	"github.com/gin-gonic/gin"
)

//Checks and authenticates the token in protected routes

func Auth(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")

	if authHeader == "" || len(authHeader) < 7 {
		generalHelper.SendError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, err := userHelper.ValidateToken(authHeader)

	if err != nil {
		fmt.Print(userID)
		generalHelper.SendError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	c.Next()
}
