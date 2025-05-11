package routers

import (
	"twitter/src/handlers"
	"twitter/src/middlewares"

	"github.com/gin-gonic/gin"
)

func TweetRouter(r *gin.RouterGroup) {
	h := handlers.NewTweetHelper()
	r.POST("/post", middlewares.GetIdFromUsernameAndPassword, h.PostTweet)
	r.GET("/get/tweet", h.GetTweet)
	r.GET("/get/tweets", middlewares.GetIdFromUsername, h.GetTweets)
}