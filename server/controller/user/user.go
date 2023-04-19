package controller

import (
	"fmt"
	"net/http"

	"p8ion/database"
	"p8ion/server/model"
	"p8ion/utils"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type SignupRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Image struct {
	Image64 string `json:"image64"`
}

func SignupUser(c *gin.Context) {
	var req SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	println(req.Username)
	user := model.User{
		ID:       1,
		Username: req.Username,
		Email:    req.Email,
	}

	db := database.GetDB()

	if err := db.Create(&user).Error; err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendResponse(c, http.StatusOK, "User created successfully")
}

func Dummy(c *gin.Context) {
	utils.SendResponse(c, http.StatusOK, "Hi User, This is a test route")
}

func ParseImage(c *gin.Context) {
	var image64 Image
	if err := c.ShouldBindBodyWith(&image64, binding.JSON); err != nil {
		utils.SendResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Print(image64)
	c.JSON(http.StatusAccepted, &image64)

}

// func GetUserFromJwt(c *gin.Context) {
// 	authHeader := c.Request.Header.Get("Authorization")

// 	userId, err := authHelper.ValidateToken(authHeader)

// 	if err != nil {
// 		utils.SendResponse(c, http.StatusUnauthorized, "Unauthorised")
// 		return
// 	}

// 	var user User

// 	db := database.GetDB()
// 	err = db.First(&user, userId).Error

// 	if err != nil {
// 		utils.SendResponse(c, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	generalHelper.SendResponse(c, http.StatusOK, user)
// 	return
// }

func GetImageData(c *gin.Context) {
	userId := c.Param("userId")
	println(userId)
	utils.SendResponse(c, http.StatusAccepted, "Images")

}
