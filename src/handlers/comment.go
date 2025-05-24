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

func (h *CommentHelper) UpdateComment(ctx *gin.Context) {
	req := dtos.CommentUpdate{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithValidationError(http.StatusNotAcceptable, err, "validation error"))
		return
	}
	res, err := h.Service.Update(ctx, req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "error in updating comment"))
		return
	}
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, res, "comment updated successfuly"))
}

func (h *CommentHelper) DeleteComment(ctx *gin.Context) {
	err := h.Service.Delete(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "error in delete comment"))
		return
	}
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, nil, "comment deleted successfuly"))
}

func (h *CommentHelper) GetComment(ctx *gin.Context) {
	comment_id := ctx.Query("comment_id")
	if len(comment_id) == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, fmt.Errorf("error in query"), "enter comment_id"))
		return
	}
	ctx.Set("comment_id", comment_id)
	res, err := h.Service.GetCommentById(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "error in get comment"))
		return
	}
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, res, "comment recived"))
}