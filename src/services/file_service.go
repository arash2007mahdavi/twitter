package services

import (
	"context"
	"encoding/base64"
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
	file, _ := TypeComverter[models.File](upload)
	user := ctx.Value("user_id").(int)
	file.CreatedBy = user
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

func (s *FileService) GetFileById(ctx context.Context) (dtos.FileResponse, error) {
	file_id := ctx.Value("file_id")
	tx := s.Database.WithContext(ctx).Begin()
	file := models.File{}
	err := tx.Model(&models.File{}).Where("id = ?", file_id).First(&file).Error
	if err != nil {
		return dtos.FileResponse{}, err
	}
	file_res, _:= TypeComverter[dtos.FileResponse](file)

	data, err := os.Open(fmt.Sprintf("%s/%s", file.Directory, file.Name))
	if err != nil {
		return dtos.FileResponse{}, err
	}
	defer data.Close()

	fileBytes, err := io.ReadAll(data)
	if err != nil {
		return dtos.FileResponse{}, err
	}
	tx.Commit()

	file_binery := base64.StdEncoding.EncodeToString(fileBytes)
	file_res.Binery = file_binery
	return *file_res, nil
}

func (s *FileService) DeleteFileById(ctx context.Context) error {
	user_id := ctx.Value("user_id")
	file_id := ctx.Value("file_id")

	tx := s.Database.WithContext(ctx).Begin()
	file := models.File{}
	err := tx.Model(&models.File{}).Where("id = ?", file_id).First(&file).Error
	if err != nil {
		return err
	}
	if file.CreatedBy != user_id {
		return fmt.Errorf("user cant delete file")
	}
	file_name := fmt.Sprintf("%s/%s", file.Directory, file.Name)
	err = os.Remove(file_name)
	if err != nil {
		return err
	}
	err = tx.Model(&models.File{}).Where("id = ?", file_id).Delete(&models.File{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}