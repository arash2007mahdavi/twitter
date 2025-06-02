package handlers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"twitter/src/database"
	"twitter/src/dtos"
	"twitter/src/logger"
	"twitter/src/responses"
	"twitter/src/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FileHelper struct {
	Logger   logger.Logger
	Service  services.FileService
	Database *gorm.DB
}

func NewFileHelper() *FileHelper {
	return &FileHelper{
		Logger: logger.NewLogger(),
		Service: *services.NewFileService(),
		Database: database.GetDB(),
	}
}

// Create godoc
// @Summary Create File
// @Description create new file
// @Tags File
// @Accept x-www-form-urlencoded
// @Produce json
// @Param data formData dtos.UploadFileRequest true "new file"
// @Param file formData file true "create file"
// @Success 200 {object} responses.Response{result=[]dtos.UserResponse} "Success"
// @Failure 406 {object} responses.Response{} "Error"
// @Router /file/post [post]
func (h *FileHelper) Create(ctx *gin.Context) {
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
	req.Name, err = saveUploadFile(upload.File, req.Directory)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, responses.GenerateResponseWithError(http.StatusBadRequest, err, "error in save file"))
		return
	}

	res, err := h.Service.Create(ctx, &req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, responses.GenerateResponseWithError(http.StatusBadRequest, err, "error in add file to database"))
		return
	}
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, res, "file saved successfuly"))
}

func saveUploadFile(file *multipart.FileHeader, directory string) (string, error) {
	randFileName := uuid.New()
	err := os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		return "", err
	}
	fileName := file.Filename
	fileNameArr := strings.Split(fileName, ".")
	fileExt := fileNameArr[len(fileNameArr) - 1]
	fileName = fmt.Sprintf("%s.%s", randFileName, fileExt)
	dst := fmt.Sprintf("%s/%s", directory, fileName)

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		return "", err
	}
	return fileName, nil
}