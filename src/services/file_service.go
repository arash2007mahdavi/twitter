package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"twitter/src/database"
	"twitter/src/database/models"
	"twitter/src/dtos"
	"twitter/src/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FileService struct {
	Database *gorm.DB
	Logger   logger.Logger
}

func NewFileService() *FileService {
	return &FileService{
		Database: database.GetDB(),
		Logger:   logger.NewLogger(),
	}
}

func (s *FileService) Create(ctx context.Context, upload *dtos.CreateFileRequest) (*dtos.FileResponse, error) {
	tx := s.Database.WithContext(ctx).Begin()
	file, _:= TypeComverter[models.File](upload)
	err := tx.Create(&file).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	res, _ := TypeComverter[dtos.FileResponse](file)
	return res, nil
}

func SaveUploadFile(file *multipart.FileHeader, directory string) (string, error) {
	randFileName := uuid.New()
	err := os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		return "", err
	}
	fileName := file.Filename
	fileNameArr := strings.Split(fileName, ".")
	fileExt := fileNameArr[len(fileNameArr)-1]
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
