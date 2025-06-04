package handlers

import (
	"net/http"
	"strconv"
	"twitter/src/database"
	"twitter/src/dtos"
	"twitter/src/logger"
	"twitter/src/responses"
	"twitter/src/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FileHelper struct {
	Logger   logger.Logger
	Service  services.FileService
	Database *gorm.DB
}

func NewFileHelper() *FileHelper {
	return &FileHelper{
		Logger:   logger.NewLogger(),
		Service:  *services.NewFileService(),
		Database: database.GetDB(),
	}
}

// TweetFile godoc
// @Summary Create File For Tweet
// @Description create new file for tweet
// @Tags File
// @Accept x-www-form-urlencoded
// @Produce json
// @Param data formData dtos.UploadFileRequest true "new file-data"
// @Param file formData file true "new file"
// @Param username query string true "user's username"
// @Param password query string true "user's password"
// @Param tweet_id query int true "tweet id"
// @Success 200 {object} responses.Response{result=dtos.FileResponse} "Success"
// @Failure 400 {object} responses.Response{} "Validation Error"
// @Failure 424 {object} responses.Response{} "Error"
// @Router /file/post/tweet [post]
func (h *FileHelper) TweetFile(ctx *gin.Context) {
	upload := dtos.UploadFileRequest{}
	err := ctx.ShouldBind(&upload)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, responses.GenerateResponseWithValidationError(http.StatusNotAcceptable, err, "validation error"))
		return
	}
	req := dtos.CreateFileRequest{}
	req.Description = upload.Description
	req.MimeType = upload.File.Header.Get("Content-Type")
	req.Directory = "uploads"
	tweet_id := ctx.Value("tweet_id")
	tweet_id_int, _ := strconv.Atoi(tweet_id.(string))
	req.TweetId = &tweet_id_int
	req.Name, err = services.SaveUploadFile(upload.File, req.Directory)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusFailedDependency, responses.GenerateResponseWithError(http.StatusFailedDependency, err, "error in save file"))
		return
	}

	res, err := h.Service.Create(ctx, &req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, responses.GenerateResponseWithError(http.StatusBadRequest, err, "error in add file to database"))
		return
	}
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, res, "file saved successfuly"))
}

// CommentFile godoc
// @Summary Create File For Comment
// @Description create new file for comment
// @Tags File
// @Accept x-www-form-urlencoded
// @Produce json
// @Param data formData dtos.UploadFileRequest true "new file-data"
// @Param file formData file true "new file"
// @Param username query string true "user's username"
// @Param password query string true "user's password"
// @Param comment_id query int true "comment id"
// @Success 200 {object} responses.Response{result=dtos.FileResponse} "Success"
// @Failure 400 {object} responses.Response{} "Validation Error"
// @Failure 424 {object} responses.Response{} "Error"
// @Router /file/post/comment [post]
func (h *FileHelper) CommentFile(ctx *gin.Context) {
	upload := dtos.UploadFileRequest{}
	err := ctx.ShouldBind(&upload)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, responses.GenerateResponseWithValidationError(http.StatusNotAcceptable, err, "validation error"))
		return
	}
	req := dtos.CreateFileRequest{}
	req.Description = upload.Description
	req.MimeType = upload.File.Header.Get("Content-Type")
	req.Directory = "uploads"
	comment_id := ctx.Value("comment_id")
	comment_id_int, _ := strconv.Atoi(comment_id.(string))
	req.CommentId = &comment_id_int
	req.Name, err = services.SaveUploadFile(upload.File, req.Directory)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusFailedDependency, responses.GenerateResponseWithError(http.StatusFailedDependency, err, "error in save file"))
		return
	}

	res, err := h.Service.Create(ctx, &req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, responses.GenerateResponseWithError(http.StatusBadRequest, err, "error in add file to database"))
		return
	}
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, res, "file saved successfuly"))
}

// GetFile godoc
// @Summary Get File With Id
// @Description get file with id
// @Tags File
// @Produce json
// @Param file_id query int true "file id"
// @Success 200 {object} responses.Response{result=dtos.FileResponse} "Success"
// @Failure 500 {object} responses.Response{} "Error"
// @Router /file/get/file [get]
func (h *FileHelper) GetFile(ctx *gin.Context) {
	file_id := ctx.Query("file_id")
	ctx.Set("file_id", file_id)
	file, err := h.Service.GetFileById(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, responses.GenerateResponseWithError(http.StatusInternalServerError, err, "error in get file from database"))
		return
	}
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, file, "file got"))
}

// GetFileInformation godoc
// @Summary Get File Information With Id
// @Description get file information with id
// @Tags File
// @Produce json
// @Param file_id query int true "file id"
// @Success 200 {object} responses.Response{result=models.File} "Success"
// @Failure 500 {object} responses.Response{} "Error"
// @Router /file/get/information [get]
func (h *FileHelper) GetFileInformation(ctx *gin.Context) {
	file_id := ctx.Query("file_id")
	ctx.Set("file_id", file_id)
	file, err := h.Service.GetFileInformationById(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, responses.GenerateResponseWithError(http.StatusInternalServerError, err, "error in get file information from database"))
		return
	}
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, file, "file got"))
}

// DeleteFile godoc
// @Summary Delete File With Id
// @Description delete file with id and check its owner
// @Tags File
// @Produce json
// @Param file_id query int true "file id"
// @Param username query string true "user's username"
// @Param password query string true "user's password"
// @Success 200 {object} responses.Response{result=string} "Success"
// @Failure 500 {object} responses.Response{} "Error"
// @Router /file/delete [delete]
func (h *FileHelper) DeleteFile(ctx *gin.Context) {
	file_id := ctx.Query("file_id")
	ctx.Set("file_id", file_id)
	err := h.Service.DeleteFileById(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, responses.GenerateResponseWithError(http.StatusInternalServerError, err, "error in delete file from database"))
		return
	}
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, "deleted", "file deleted successfuly"))
}