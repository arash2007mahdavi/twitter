package services

import (
	"twitter/src/database"
	"twitter/src/database/models"
	"twitter/src/dtos"
)

type UserService struct {
	Base *database.BaseService[models.User, dtos.UserCreate, dtos.UserUpdate, dtos.UserResponse]
}

func NewUserService() *UserService {
	return &UserService{Base: database.NewBaseService[models.User, dtos.UserCreate, dtos.UserUpdate, dtos.UserResponse]()}
}