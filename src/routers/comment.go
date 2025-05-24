package routers

import (
	"twitter/src/handlers"
	"twitter/src/middlewares"

	"github.com/gin-gonic/gin"
)

func CommentRouter(r *gin.RouterGroup) {
	h := handlers.NewCommentHelper()
	r.POST("/post", middlewares.GetIdForComment, h.PostComment)
	r.PUT("/update", middlewares.CheckCommentOwner, h.UpdateComment)
	r.DELETE("/delete", middlewares.CheckForDeleteComment, h.DeleteComment)
	r.GET("/get/comment", h.GetComment)
	r.GET("/get/comments", middlewares.GetIdFromUsername, h.GetComments)
}