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

type CommentHelper struct {
	Logger   logger.Logger
	Service  *services.CommentService
	Database *gorm.DB
}

func NewCommentHelper() *CommentHelper {
	return &CommentHelper{
		Logger:   logger.NewLogger(),
		Service:  services.NewCommentService(),
		Database: database.GetDB(),
	}
}

func (h *CommentHelper) PostComment(ctx *gin.Context) {
	req := dtos.CommentCreate{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithValidationError(http.StatusNotAcceptable, err, "validation error"))
		return
	}
	res, err := h.Service.PostComment(ctx, &req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "error in post comment"))
		return
	}
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, res, "comment posted successfuly"))
}