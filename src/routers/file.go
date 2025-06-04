package routers

import (
	"twitter/src/handlers"
	"twitter/src/middlewares"

	"github.com/gin-gonic/gin"
)

func FileRouter(r *gin.RouterGroup) {
	h := handlers.NewFileHelper()
	r.GET("/get/file", h.GetFile)
	r.GET("/get/information", h.GetFileInformation)
	r.DELETE("/delete", middlewares.GetIdFromUsernameAndPassword, h.DeleteFile)
	r.POST("/post/tweet", middlewares.CheckTweetForAddFile, h.TweetFile)
	r.POST("/post/comment", middlewares.CheckCommentForAddFile, h.CommentFile)
}
