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
	h.Logger.Info(logger.Tweet, logger.New, "new tweet posted", map[logger.ExtraCategory]interface{}{logger.Userid: res.UserId, logger.Tweetid: res.Id})
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, res, "new tweet posted"))
}

func (h *TweetHelper) GetTweet(ctx *gin.Context) {
	tweet_id := ctx.Query("tweet_id")
	if tweet_id == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, responses.GenerateResponseWithError(http.StatusBadRequest, nil, "tweet id is required"))
		return
	}
	ctx.Set("tweet_id", tweet_id)
	tweet, err := h.Service.GetTweetByID(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, responses.GenerateResponseWithError(http.StatusNotFound, err, "error in get tweet"))
		return
	}
	h.Logger.Info(logger.Tweet, logger.Get, "tweet got", map[logger.ExtraCategory]interface{}{logger.Tweetid: tweet_id})
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, tweet, "tweet retrieved"))
}

func (h *TweetHelper) GetTweets(ctx *gin.Context) {
	tweets, err := h.Service.GetTweets(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "error in get tweets"))
		return
	}
	id := ctx.Value("user_id")
	h.Logger.Info(logger.Tweet, logger.Get, "tweets got", map[logger.ExtraCategory]interface{}{logger.Userid: id})
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, tweets, "users tweets"))
}