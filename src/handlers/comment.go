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

// PostComment godoc
// @Summary Post Comment
// @Description post comment with message
// @Tags Comment
// @Accept json
// @Produce json
// @Param username query string true "user's username"
// @Param password query string true "user's password"
// @Param tweet_id query string true "tweet's id"
// @Param CommentCreate body dtos.CommentCreate true "get message of new comment"
// @Success 200 {object} responses.Response{result=dtos.CommentResponse} "Success"
// @Failure 400 {object} responses.Response{} "Validation Error"
// @Failure 500 {object} responses.Response{} "Internal Server Error"
// @Router /comment/post [post]
func (h *CommentHelper) PostComment(ctx *gin.Context) {
	req := dtos.CommentCreate{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, responses.GenerateResponseWithValidationError(http.StatusBadRequest, err, "validation error"))
		return
	}
	res, err := h.Service.PostComment(ctx, &req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, responses.GenerateResponseWithError(http.StatusInternalServerError, err, "error in post comment"))
		return
	}
	h.Logger.Info(logger.Comment, logger.Add, "new comment posted", map[logger.ExtraCategory]interface{}{logger.Commentid: res.Id})
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, res, "comment posted successfuly"))
}

// UpdateComment godoc
// @Summary Update Comment
// @Description update comment's message
// @Tags Comment
// @Accept json
// @Produce json
// @Param username query string true "user's username"
// @Param password query string true "user's password"
// @Param comment_id query string true "comment's id"
// @Param CommentUpdate body dtos.CommentUpdate true "get updated message of comment"
// @Success 200 {object} responses.Response{result=dtos.CommentResponse} "Success"
// @Failure 400 {object} responses.Response{} "Validation Error"
// @Failure 500 {object} responses.Response{} "Internal Server Error"
// @Router /comment/update [put]
func (h *CommentHelper) UpdateComment(ctx *gin.Context) {
	req := dtos.CommentUpdate{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, responses.GenerateResponseWithValidationError(http.StatusBadRequest, err, "validation error"))
		return
	}
	res, err := h.Service.Update(ctx, req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, responses.GenerateResponseWithError(http.StatusInternalServerError, err, "error in updating comment"))
		return
	}
	h.Logger.Info(logger.Comment, logger.Update, "comment updated", map[logger.ExtraCategory]interface{}{logger.Commentid: ctx.Value("comment_id")})
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, res, "comment updated successfuly"))
}

// DeleteComment godoc
// @Summary Delete Comment
// @Description delete comment
// @Tags Comment
// @Produce json
// @Param username query string true "user's username"
// @Param password query string true "user's password"
// @Param comment_id query string true "comment's id"
// @Success 200 {object} responses.Response{result=string} "Success"
// @Failure 500 {object} responses.Response{} "Internal Server Error"
// @Router /comment/delete [delete]
func (h *CommentHelper) DeleteComment(ctx *gin.Context) {
	err := h.Service.Delete(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, responses.GenerateResponseWithError(http.StatusInternalServerError, err, "error in delete comment"))
		return
	}
	h.Logger.Info(logger.Comment, logger.Delete, "comment deleted", map[logger.ExtraCategory]interface{}{logger.Commentid: ctx.Value("comment_id")})
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, "deleted", "comment deleted successfuly"))
}

// GetComment godoc
// @Summary Get Comment
// @Description get comment
// @Tags Comment
// @Produce json
// @Param comment_id query string true "comment's id"
// @Success 200 {object} responses.Response{result=models.Comment} "Success"
// @Failure 500 {object} responses.Response{} "Internal Server Error"
// @Router /comment/get/comment [get]
func (h *CommentHelper) GetComment(ctx *gin.Context) {
	comment_id := ctx.Query("comment_id")
	if len(comment_id) == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, fmt.Errorf("error in query"), "enter comment_id"))
		return
	}
	ctx.Set("comment_id", comment_id)
	res, err := h.Service.GetCommentById(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, responses.GenerateResponseWithError(http.StatusInternalServerError, err, "error in get comment"))
		return
	}
	h.Logger.Info(logger.Comment, logger.Get, "comment get", map[logger.ExtraCategory]interface{}{logger.Commentid: comment_id})
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, res, "comment recived"))
}

// GetComments godoc
// @Summary Get Comments
// @Description get comments of an user
// @Tags Comment
// @Produce json
// @Param username query string true "user's username"
// @Success 200 {object} responses.Response{result=[]dtos.CommentResponse} "Success"
// @Failure 500 {object} responses.Response{} "Internal Server Error"
// @Router /comment/get/comments [get]
func (h *CommentHelper) GetComments(ctx *gin.Context) {
	res, err := h.Service.GetComments(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "error in get comments"))
		return
	}
	h.Logger.Info(logger.Comment, logger.Get, "comments get", map[logger.ExtraCategory]interface{}{logger.Userid: ctx.Value("user_id")})
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, res, "comments recived"))
}

// LikeComment godoc
// @Summary Like Comment
// @Description like comment
// @Tags Comment
// @Produce json
// @Param username query string true "user's username"
// @Param password query string true "user's password"
// @Param comment_id query string true "comment's id"
// @Success 200 {object} responses.Response{result=string} "Success"
// @Failure 500 {object} responses.Response{} "Internal Server Error"
// @Router /comment/like [post]
func (h *CommentHelper) LikeComment(ctx *gin.Context) {
	comment_id := ctx.Query("comment_id")
	if len(comment_id) == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, fmt.Errorf("error in get query"), "enter comment_id"))
		return
	}
	ctx.Set("comment_id", comment_id)
	err := h.Service.LikeComment(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, responses.GenerateResponseWithError(http.StatusInternalServerError, err, "error in like comment"))
		return
	}
	h.Logger.Info(logger.Comment, logger.Like, "comment liked", map[logger.ExtraCategory]interface{}{logger.Commentid: comment_id, logger.Userid: ctx.Value("user_id")})
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, "the comment liked successfuly", "comment liked"))
}

// DislikeComment godoc
// @Summary Dislike Comment
// @Description dislike comment
// @Tags Comment
// @Produce json
// @Param username query string true "user's username"
// @Param password query string true "user's password"
// @Param comment_id query string true "comment's id"
// @Success 200 {object} responses.Response{result=string} "Success"
// @Failure 500 {object} responses.Response{} "Internal Server Error"
// @Router /comment/dislike [post]
func (h *CommentHelper) DislikeComment(ctx *gin.Context) {
	comment_id := ctx.Query("comment_id")
	if len(comment_id) == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, fmt.Errorf("error in get query"), "enter comment_id"))
		return
	}
	ctx.Set("comment_id", comment_id)
	err := h.Service.DislikeComment(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, responses.GenerateResponseWithError(http.StatusInternalServerError, err, "error in dislike comment"))
		return
	}
	h.Logger.Info(logger.Comment, logger.Dislike, "comment disliked", map[logger.ExtraCategory]interface{}{logger.Commentid: comment_id, logger.Userid: ctx.Value("user_id")})
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, "the comment disliked successfuly", "comment disliked"))
}