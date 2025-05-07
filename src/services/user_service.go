package services

import (
	"context"
	"twitter/src/database"
	"twitter/src/database/models"
	"twitter/src/dtos"
	"twitter/src/logger"
)

type UserService struct {
	Base *database.BaseService[models.User, dtos.UserCreate, dtos.UserUpdate, dtos.UserResponse]
	Logger logger.Logger
}

func NewUserService() *UserService {
	return &UserService{Base: database.NewBaseService[models.User, dtos.UserCreate, dtos.UserUpdate, dtos.UserResponse](), Logger: logger.NewLogger()}
}

func (s *UserService) Create(ctx context.Context, req *dtos.UserCreate) (*dtos.UserResponse, error) {
	res, err := s.Base.Create(ctx, req)
	return res, err
}