package handlers

import (
	"fmt"
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
		ctx.AbortWithStatusJSON(http.StatusBadRequest, responses.GenerateResponseWithError(http.StatusBadRequest, fmt.Errorf("error in query"), "tweet id is required"))
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

func (h *TweetHelper) UpdateTweet(ctx *gin.Context) {
	req := dtos.TweetUpdate{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithValidationError(http.StatusNotAcceptable, err, "validation error"))
		return
	}
	res, err := h.Service.Update(ctx, &req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, responses.GenerateResponseWithError(http.StatusNotFound, err, "updating tweet failed"))
		return
	}
	h.Logger.Info(logger.Tweet, logger.Update, "tweet updated", map[logger.ExtraCategory]interface{}{logger.Tweetid: ctx.Value("tweet_id")})
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, res, "tweet updated successfuly"))
}

func (h *TweetHelper) DeleteTweet(ctx *gin.Context) {
	err := h.Service.Delete(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, responses.GenerateResponseWithError(http.StatusNotFound, err, "deleting tweet failed"))
		return
	}
	h.Logger.Info(logger.Tweet, logger.Delete, "tweet deleted", map[logger.ExtraCategory]interface{}{logger.Tweetid: ctx.Value("tweet_id")})
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, nil, "tweet deleted successfuly"))
}

func (h *TweetHelper) GetFollowingsTweets(ctx *gin.Context) {
	tweets, err := h.Service.GetFollowingsTweets(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, responses.GenerateResponseWithError(http.StatusNotFound, err, "error in get followers tweets"))
		return
	}
	h.Logger.Info(logger.Tweet, logger.Get, "followings tweets get", map[logger.ExtraCategory]interface{}{logger.Userid: ctx.Value("user_id")})
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, tweets, "followers tweets show"))
}

func (h *TweetHelper) TweetExplore(ctx *gin.Context) {
	res, err := h.Service.TweetExplore(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "error in load explore"))
		return
	}
	h.Logger.Info(logger.Tweet, logger.Get, "explore get", nil)
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, res, "explore loaded"))
}

func (h *TweetHelper) LikeTweet(ctx *gin.Context) {
	tweet_id := ctx.Query("tweet_id")
	if len(tweet_id) == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, fmt.Errorf("error in get query"), "enter tweet_id"))
		return
	}
	ctx.Set("tweet_id", tweet_id)
	err := h.Service.LikeTweet(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "error in like tweet"))
		return
	}
	h.Logger.Info(logger.Tweet, logger.Like, "tweet liked", map[logger.ExtraCategory]interface{}{logger.Userid: ctx.Value("user_id"), logger.Tweetid: tweet_id})
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, "liked successfuly", "tweet liked"))
}

func (h *TweetHelper) DislikeTweet(ctx *gin.Context) {
	tweet_id := ctx.Query("tweet_id")
	if len(tweet_id) == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, fmt.Errorf("error in get query"), "enter tweet_id"))
		return
	}
	ctx.Set("tweet_id", tweet_id)
	err := h.Service.DislikeTweet(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "error in dislike tweet"))
		return
	}
	h.Logger.Info(logger.Tweet, logger.Dislike, "tweet disliked", map[logger.ExtraCategory]interface{}{logger.Userid: ctx.Value("user_id"), logger.Tweetid: tweet_id})
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, "disliked successfuly", "tweet disliked"))
}