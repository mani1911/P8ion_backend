package controller

import (
	"fmt"
	"net/http"
	"p8ion/config"
	"p8ion/database"
	authHelper "p8ion/server/helpers/auth"
	helper "p8ion/server/helpers/general"
	"p8ion/server/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func OAuthRequest(c *gin.Context) {

	frontendUrl := config.GetConfig().FrontendURL

	// get the code from qs
	code := c.Request.URL.Query().Get("code")
	if len(code) == 0 {
		helper.SendError(c, http.StatusBadRequest, "Invalid Request")
		return
	}

	// get id and access token

	token, err := authHelper.GetGoogleOAuthTokens(code)

	if err != nil {
		fmt.Print("Error Parsing Token", err)
		helper.SendError(c, http.StatusInternalServerError, "Some error occurred, Refresh and Try Again")
		return
	}

	// get user with token

	user, err := authHelper.GetOAuth2User(token.AccessToken, token.IDToken)

	if err != nil {
		fmt.Print("Error Getting User", err)
		helper.SendError(c, http.StatusInternalServerError, "Some error occurred, Refresh and Try Again")
		return
	}

	Name := user.Name
	Email := user.Email

	if len(Name) == 0 || len(Email) == 0 {
		helper.SendError(c, http.StatusInternalServerError, "Unable to find User Details")
		return
	}

	// upsert user

	db := database.GetDB()
	var userDetails model.User

	if err := db.Where("Email = ?", Email).First(&userDetails).Error; err != nil {
		if err == gorm.ErrRecordNotFound {

			userReg := model.User{
				Email:    Email,
				Username: Name,
			}

			if err := db.Create(&userReg).Error; err != nil {
				helper.SendError(c, http.StatusInternalServerError, "Unable to create user, Try Again")
				return
			}

			jwtToken, err := authHelper.GenerateToken(userReg.ID)

			if err != nil {
				helper.SendError(c, http.StatusInternalServerError, "Unknown Error occurred, Try Again")
				return
			}

			helper.SendResponse(c, http.StatusOK, jwtToken)
			return
		}
		helper.SendError(c, http.StatusInternalServerError, "Unknown Error occurred, Try Again")
		return
	}

	// create jwt

	jwtToken, err := authHelper.GenerateToken(userDetails.ID)
	if err != nil {
		fmt.Print("Token Not generated:", err)
		helper.SendError(c, http.StatusInternalServerError, "Unknown Error occurred, Try Again")
		return
	}

	// set cookies

	c.SetCookie("jwt", jwtToken, 86400, "/", c.Request.URL.Hostname(), false, true)
	c.SetCookie("username", Name, 86400, "/", c.Request.URL.Hostname(), false, true)

	// redirect back to client
	fmt.Print("CLient : ", frontendUrl)
	c.Redirect(http.StatusPermanentRedirect, frontendUrl)
}
