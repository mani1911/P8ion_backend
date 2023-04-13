package controller

import (
	"net/http"

	"p8ion/database"
	"p8ion/server/model"
	"p8ion/utils"

	"github.com/gin-gonic/gin"
)

type SignupRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json: "email" binding:"required"`
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
