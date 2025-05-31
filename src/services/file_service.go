package services

import (
	"context"
	"twitter/src/database"
	"twitter/src/database/models"
	"twitter/src/dtos"
	"twitter/src/logger"

	"gorm.io/gorm"
)

type FileService struct {
	Database *gorm.DB
	Logger   logger.Logger
}

func NewFileService() *FileService {
	return &FileService{
		Database: database.GetDB(),
		Logger: logger.NewLogger(),
	}
}

func (s *FileService) Create(ctx context.Context, req *dtos.CreateFileRequest) (*dtos.FileResponse, error) {
	file, _:= TypeComverter[models.File](req)
	tx := s.Database.WithContext(ctx).Begin()
	err := tx.Create(&file).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	res, _:= TypeComverter[dtos.FileResponse](file)
	return res, nil
}