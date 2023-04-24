package controller

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"p8ion/database"
	"p8ion/server/model"
	"p8ion/utils"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sashabaranov/go-openai"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/vision/v1"
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
	Image64 string   `json:"image64"`
	Content []string `json:"content"`
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

	//Cloud Vision API
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, vision.CloudPlatformScope)
	if err != nil {
		utils.SendResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	service, err := vision.New(client)
	if err != nil {
		utils.SendResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	request := &vision.AnnotateImageRequest{
		Image: &vision.Image{
			Content: image64.Image64,
		},
		Features: []*vision.Feature{
			{
				Type: "TEXT_DETECTION",
			},
		},
	}

	batch := &vision.BatchAnnotateImagesRequest{
		Requests: []*vision.AnnotateImageRequest{request},
	}

	res, err := service.Images.Annotate(batch).Do()
	if err != nil {
		utils.SendResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var drugs []string
	var drugsWithCommas string

	if annotations := res.Responses[0].TextAnnotations; len(annotations) > 0 {
		text := annotations[0].Description
		// utils.SendResponse(c, http.StatusOK, text)
		start := 0
		for i := 0; i < len(text); i++ {
			if (text[i] == '-') || (text[i] >= '0' && text[i] <= '9') {
				break
			}
			if text[i] == '\n' {
				drugs = append(drugs, text[start:i])
				start = i + 1
			}
		}
		// utils.SendResponse(c, http.StatusOK, drugs)
		for i := 0; i < len(drugs)-1; i++ {
			drugsWithCommas += drugs[i] + ","
		}
		drugsWithCommas += drugs[len(drugs)-1]
	} else {
		utils.SendResponse(c, http.StatusBadRequest, "No Text found")
		return
	}

	//GPT API
	gpt := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	resp, err := gpt.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Describe the following drugs " + drugsWithCommas + " in 100 words with common sideeffects.",
				},
			},
		},
	)
	if err != nil {
		utils.SendResponse(c, http.StatusBadRequest, "Could not connect to GPT client")
	}

	// utils.SendResponse(c, http.StatusOK, resp.Choices[0].Message.Content)
	drugDesc := strings.Split(resp.Choices[0].Message.Content, "\n\n")
	utils.SendResponse(c, http.StatusOK, drugDesc)
}

// func putIntoDB(content string, ID uint, c *gin.Context) {
// 	db := database.GetDB()
// 	var drugDetails model.Image

// 	if err := db.Where("user_id = ?", ID).First(&drugDetails).Error; err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			drugRec := model.Image{
// 				ImageBase64: "<image-base-64>",
// 				UserID:      ID,
// 				Content:     content,
// 			}
// 			if err := db.Create(&drugRec).Error; err != nil {
// 				utils.SendResponse(c, http.StatusInternalServerError, "Error in creating Drug record")
// 				return
// 			}
// 		}
// 	}
// }

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
