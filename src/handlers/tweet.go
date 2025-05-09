package handlers

import (
	"net/http"
	"twitter/src/database"
	"twitter/src/dtos"
	"twitter/src/logger"
	"twitter/src/responses"
	"twitter/src/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TweetHelper struct {
	Logger logger.Logger
	Service *services.TweetService
	Database *gorm.DB
}

func NewTweetHelper() *TweetHelper {
	return &TweetHelper{
		Logger: logger.NewLogger(),
		Service: services.NewTweetService(),
		Database: database.GetDB(),
	}
}

func (h *TweetHelper) PostTweet(ctx *gin.Context) {
	req := dtos.TweetCreate{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithValidationError(http.StatusNotAcceptable, err, "validation error"))
		return
	}
	res, err := h.Service.NewTweet(ctx, &req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "error in post tweet"))
		return
	}
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, res, "new tweet posted"))
}