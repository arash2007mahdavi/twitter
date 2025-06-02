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

// PostTweet godoc
// @Summary Post Tweet
// @Description tweet post
// @Tags Tweet
// @Accept json
// @Produce json
// @Param username query string true "get user's username"
// @Param password query string true "get user's password"
// @Param TweetCreate body dtos.TweetCreate true "get title and message of new tweet"
// @Success 200 {object} responses.Response{result=dtos.TweetResponse} "Success"
// @Failure 400 {object} responses.Response{} "Validation Error"
// @Failure 500 {object} responses.Response{} "Internal Server Error"
// @Router /tweet/post [post]
func (h *TweetHelper) PostTweet(ctx *gin.Context) {
	req := dtos.TweetCreate{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, responses.GenerateResponseWithValidationError(http.StatusBadRequest, err, "validation error"))
		return
	}
	res, err := h.Service.NewTweet(ctx, &req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, responses.GenerateResponseWithError(http.StatusInternalServerError, err, "error in post tweet"))
		return
	}
	h.Logger.Info(logger.Tweet, logger.New, "new tweet posted", map[logger.ExtraCategory]interface{}{logger.Userid: res.UserId, logger.Tweetid: res.Id})
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, res, "new tweet posted"))
}

// GetTweet godoc
// @Summary Get Tweet
// @Description get tweet by id
// @Tags Tweet
// @Produce json
// @Param tweet_id query string true "tweet's id"
// @Success 200 {object} responses.Response{result=models.Tweet} "Success"
// @Failure 400 {object} responses.Response{} "Validation Error"
// @Failure 500 {object} responses.Response{} "Internal Server Error"
// @Router /tweet/get/tweet [get]
func (h *TweetHelper) GetTweet(ctx *gin.Context) {
	tweet_id := ctx.Query("tweet_id")
	if tweet_id == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, responses.GenerateResponseWithError(http.StatusBadRequest, fmt.Errorf("error in query"), "tweet id is required"))
		return
	}
	ctx.Set("tweet_id", tweet_id)
	tweet, err := h.Service.GetTweetByID(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, responses.GenerateResponseWithError(http.StatusInternalServerError, err, "error in get tweet"))
		return
	}
	h.Logger.Info(logger.Tweet, logger.Get, "tweet got", map[logger.ExtraCategory]interface{}{logger.Tweetid: tweet_id})
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, tweet, "tweet retrieved"))
}

// GetTweets godoc
// @Summary Get Tweets
// @Description get tweets those tweeted by an user
// @Tags Tweet
// @Produce json
// @Param username query string true "get user's username"
// @Success 200 {object} responses.Response{result=[]dtos.TweetResponse} "Success"
// @Failure 500 {object} responses.Response{} "Internal Server Error"
// @Router /tweet/get/tweets [get]
func (h *TweetHelper) GetTweets(ctx *gin.Context) {
	tweets, err := h.Service.GetTweets(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, responses.GenerateResponseWithError(http.StatusInternalServerError, err, "error in get tweets"))
		return
	}
	id := ctx.Value("user_id")
	h.Logger.Info(logger.Tweet, logger.Get, "tweets got", map[logger.ExtraCategory]interface{}{logger.Userid: id})
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, tweets, "users tweets"))
}

// UpdateTweet godoc
// @Summary Update Tweet
// @Description update tweet's title and message
// @Tags Tweet
// @Accept json
// @Produce json
// @Param username query string true "owner's username"
// @Param password query string true "owner's password"
// @Param tweet_id query string true "tweet's id"
// @Param TweetUpdate body dtos.TweetUpdate true "update's fields"
// @Success 200 {object} responses.Response{result=dtos.TweetResponse} "Success"
// @Failure 400 {object} responses.Response{} "Validation Error"
// @Failure 500 {object} responses.Response{} "Internal Server Error"
// @Router /tweet/update [put]
func (h *TweetHelper) UpdateTweet(ctx *gin.Context) {
	req := dtos.TweetUpdate{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, responses.GenerateResponseWithValidationError(http.StatusBadRequest, err, "validation error"))
		return
	}
	res, err := h.Service.Update(ctx, &req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, responses.GenerateResponseWithError(http.StatusInternalServerError, err, "updating tweet failed"))
		return
	}
	h.Logger.Info(logger.Tweet, logger.Update, "tweet updated", map[logger.ExtraCategory]interface{}{logger.Tweetid: ctx.Value("tweet_id")})
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, res, "tweet updated successfuly"))
}

// DeleteTweet godoc
// @Summary Delete Tweet
// @Description delete tweet by id and its owner information
// @Tags Tweet
// @Produce json
// @Param username query string true "owner's username"
// @Param password query string true "owner's password"
// @Param tweet_id query string true "tweet's id"
// @Success 200 {object} responses.Response{result=string} "Success"
// @Failure 500 {object} responses.Response{} "Internal Server Error"
// @Router /tweet/delete [delete]
func (h *TweetHelper) DeleteTweet(ctx *gin.Context) {
	err := h.Service.Delete(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, responses.GenerateResponseWithError(http.StatusInternalServerError, err, "deleting tweet failed"))
		return
	}
	h.Logger.Info(logger.Tweet, logger.Delete, "tweet deleted", map[logger.ExtraCategory]interface{}{logger.Tweetid: ctx.Value("tweet_id")})
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, "deleted", "tweet deleted successfuly"))
}

// GetFollowingsTweets godoc
// @Summary Get Followings Tweets
// @Description delete tweet by id and its owner information
// @Tags Tweet
// @Produce json
// @Param username query string true "owner's username"
// @Param password query string true "owner's password"
// @Success 200 {object} responses.Response{result=[]dtos.TweetResponse} "Success"
// @Failure 500 {object} responses.Response{} "Internal Server Error"
// @Router /tweet/get/followings/tweets [get]
func (h *TweetHelper) GetFollowingsTweets(ctx *gin.Context) {
	tweets, err := h.Service.GetFollowingsTweets(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, responses.GenerateResponseWithError(http.StatusInternalServerError, err, "error in get followers tweets"))
		return
	}
	h.Logger.Info(logger.Tweet, logger.Get, "followings tweets get", map[logger.ExtraCategory]interface{}{logger.Userid: ctx.Value("user_id")})
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, tweets, "followers tweets show"))
}

// GetFollowingsTweets godoc
// @Summary Get Followings Tweets
// @Description delete tweet by id and its owner information
// @Tags Tweet
// @Produce json
// @Param username query string true "owner's username"
// @Param password query string true "owner's password"
// @Success 200 {object} responses.Response{result=[]dtos.TweetResponse} "Success"
// @Failure 500 {object} responses.Response{} "Internal Server Error"
// @Router /tweet/explore [get]
func (h *TweetHelper) TweetExplore(ctx *gin.Context) {
	res, err := h.Service.TweetExplore(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "error in load explore"))
		return
	}
	h.Logger.Info(logger.Tweet, logger.Get, "explore get", nil)
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, res, "explore loaded"))
}

// LikeTweet godoc
// @Summary Like Tweet
// @Description like tweet
// @Tags Tweet
// @Produce json
// @Param username query string true "owner's username"
// @Param password query string true "owner's password"
// @Param tweet_id query string true "tweet's id"
// @Success 200 {object} responses.Response{result=string} "Success"
// @Failure 500 {object} responses.Response{} "Internal Server Error"
// @Router /tweet/like [post]
func (h *TweetHelper) LikeTweet(ctx *gin.Context) {
	tweet_id := ctx.Query("tweet_id")
	if len(tweet_id) == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, fmt.Errorf("error in get query"), "enter tweet_id"))
		return
	}
	ctx.Set("tweet_id", tweet_id)
	err := h.Service.LikeTweet(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, responses.GenerateResponseWithError(http.StatusInternalServerError, err, "error in like tweet"))
		return
	}
	h.Logger.Info(logger.Tweet, logger.Like, "tweet liked", map[logger.ExtraCategory]interface{}{logger.Userid: ctx.Value("user_id"), logger.Tweetid: tweet_id})
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, "liked successfuly", "tweet liked"))
}

// DislikeTweet godoc
// @Summary Dislike Tweet
// @Description dislike tweet
// @Tags Tweet
// @Produce json
// @Param username query string true "owner's username"
// @Param password query string true "owner's password"
// @Param tweet_id query string true "tweet's id"
// @Success 200 {object} responses.Response{result=string} "Success"
// @Failure 500 {object} responses.Response{} "Internal Server Error"
// @Router /tweet/dislike [post]
func (h *TweetHelper) DislikeTweet(ctx *gin.Context) {
	tweet_id := ctx.Query("tweet_id")
	if len(tweet_id) == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, fmt.Errorf("error in get query"), "enter tweet_id"))
		return
	}
	ctx.Set("tweet_id", tweet_id)
	err := h.Service.DislikeTweet(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, responses.GenerateResponseWithError(http.StatusInternalServerError, err, "error in dislike tweet"))
		return
	}
	h.Logger.Info(logger.Tweet, logger.Dislike, "tweet disliked", map[logger.ExtraCategory]interface{}{logger.Userid: ctx.Value("user_id"), logger.Tweetid: tweet_id})
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, "disliked successfuly", "tweet disliked"))
}