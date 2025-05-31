package routers

import (
	"twitter/src/handlers"

	"github.com/gin-gonic/gin"
)

func FileRouter(r *gin.RouterGroup) {
	h := handlers.NewFileHelper()
	r.POST("/post", h.Create)
}