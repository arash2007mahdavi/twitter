package routers

import (
	"twitter/src/handlers"
	"twitter/src/middlewares"

	"github.com/gin-gonic/gin"
)

func FileRouter(r *gin.RouterGroup) {
	h := handlers.NewFileHelper()
	r.POST("/post/tweet", middlewares.CheckTweetForAddFile, h.TweetFile)
}
