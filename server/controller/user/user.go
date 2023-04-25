package controller

import (
	"context"
	"net/http"
	"os"

	"p8ion/database"
	authHelper "p8ion/server/helpers/auth"
	helper "p8ion/server/helpers/general"
	"p8ion/server/model"
	"p8ion/utils"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sashabaranov/go-openai"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/vision/v1"
)

type User struct {
	ID       uint `gorm:"primarykey"`
	Username string
	Email    string
}

type Image struct {
	Image64 string `json:"image64"`
}

func ParseImage(c *gin.Context) {
	db := database.GetDB()

	var image64 Image
	if err := c.ShouldBindBodyWith(&image64, binding.JSON); err != nil {
		utils.SendResponse(c, http.StatusBadRequest, err.Error())
		return
	}

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
					Content: "Describe the following drugs " + drugsWithCommas + " in 5 words with common sideeffects.",
				},
			},
		},
	)
	if err != nil {
		utils.SendResponse(c, http.StatusBadRequest, "Could not connect to GPT client")
		return
	}

	// utils.SendResponse(c, http.StatusOK, resp.Choices[0].Message.Content)
	//drugDesc := strings.Split(resp.Choices[0].Message.Content, "\n\n")
	userID, err := authHelper.ValidateToken(c.Request.Header.Get("Authorization"))

	//print(drugDesc)
	print("User ID : ", userID)

	var user User
	result := db.First(&user, userID)
	if result.Error != nil {
		utils.SendResponse(c, http.StatusNotFound, result.Error.Error())
		return
	}

	imageData := model.Image{
		UserID:      userID,
		ImageBase64: image64.Image64,
		Content:     resp.Choices[0].Message.Content,
	}

	if err := db.Create(&imageData).Error; err != nil {
		helper.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendResponse(c, http.StatusOK, "Image Saved Successfully")
}

func GetImageData(c *gin.Context) {
	db := database.GetDB()
	userID, err := authHelper.ValidateToken(c.Request.Header.Get("Authorization"))

	print(userID)
	if err != nil {
		helper.SendError(c, http.StatusUnauthorized, "Unauthorised")
	}
	var images []model.Image

	result := db.Find(&images, "user_id = ?", userID)
	if result.Error != nil {
		helper.SendError(c, http.StatusNotFound, result.Error.Error())
		return
	}

	print(images)

	utils.SendResponse(c, http.StatusAccepted, images)
}
