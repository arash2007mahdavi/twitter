package services

import (
	"context"
	"twitter/src/database"
	"twitter/src/database/models"
	"twitter/src/dtos"
	"twitter/src/logger"

	"gorm.io/gorm"
)

type UserService struct {
	Base *database.BaseService[models.User, dtos.UserCreate, dtos.UserUpdate, dtos.UserResponse]
	Logger logger.Logger
	DB *gorm.DB
}

func NewUserService() *UserService {
	return &UserService{Base: database.NewBaseService[models.User, dtos.UserCreate, dtos.UserUpdate, dtos.UserResponse](),
		Logger: logger.NewLogger(),
		DB: database.GetDB(),
	}
}

func (s *UserService) Create(ctx context.Context, req *dtos.UserCreate) (*dtos.UserResponse, error) {
	res, err := s.Base.Create(ctx, req)
	return res, err
}

func (s *UserService) Update(ctx context.Context, req *dtos.UserUpdate) (*dtos.UserResponse, error) {
	res, err := s.Base.Update(ctx, req)
	return res, err
}

func (s *UserService) Delete(ctx context.Context) error {
	err := s.Base.Delete(ctx)
	return err
}

func (s *UserService) GetUsers(ctx context.Context) (*[]models.User, error) {
	return s.Base.GetAll(ctx)
}

func (s *UserService) GetProfile(ctx context.Context) (*models.User, error) {
	user := models.User{}
	username := ctx.Value("Username")
	password := ctx.Value("Password")
	s.DB.Model(&user).Preload("Tweets").Preload("Comments").Preload("Followers").Preload("Followings").Where("username = ? AND password = ?", username, password).First(&user)
	return &user, nil
}