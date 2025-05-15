package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"twitter/src/database"
	"twitter/src/database/models"
	"twitter/src/dtos"
	"twitter/src/logger"

	"gorm.io/gorm"
)

func TypeComverter[T any](req any) (*T, error) {
	var result T
	json_sample, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(json_sample, &result)
	if err != nil {
		return nil ,err
	}
	return &result, nil
}

type UserService struct {
	Logger logger.Logger
	DB *gorm.DB
}

func NewUserService() *UserService {
	return &UserService{
		Logger: logger.NewLogger(),
		DB: database.GetDB(),
	}
}

func (s *UserService) Create(ctx context.Context, req *dtos.UserCreate) (*dtos.UserResponse, error) {
	data, err := TypeComverter[models.User](req)
	if err != nil {
		return nil, fmt.Errorf("error in comverting model")
	}
	tx := s.DB.WithContext(ctx).Begin()
	err = tx.Create(data).Error
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error in creating new model")
	}
	tx.Commit()
	return TypeComverter[dtos.UserResponse](data)
}

func (s *UserService) Update(ctx context.Context, req *dtos.UserUpdate) (*dtos.UserResponse, error) {
	data, err := TypeComverter[map[string]interface{}](req)
	if err != nil {
		return nil, fmt.Errorf("error in comverting model")
	}
	int_m, _:= strconv.Atoi(ctx.Value("modified_by").(string))
	int_id := ctx.Value("user_id").(int)
	(*data)["modified_by"] = sql.NullInt64{Int64: int64(int_m), Valid: true}
	(*data)["modified_at"] = sql.NullTime{Time: time.Now().UTC(), Valid: true}
	tx := s.DB.WithContext(ctx).Begin()
	model := new(models.User)
	err = tx.Model(model).
		Where("id = ? AND deleted_by is null", ctx.Value("user_id")).
		Updates(*data).Error
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error in updating model")
	}
	user := models.User{}
	tx.Model(&models.User{}).Where("id = ?", int_id).First(&user)
	tx.Commit()
	return TypeComverter[dtos.UserResponse](user)
}

func (s *UserService) Delete(ctx context.Context) error {
	data := map[string]interface{}{}
	int_d, _:= strconv.Atoi(ctx.Value("deleted_by").(string))
	int_id := ctx.Value("user_id").(int)
	(data)["deleted_by"] = sql.NullInt64{Int64: int64(int_d), Valid: true}
	(data)["deleted_at"] = sql.NullTime{Time: time.Now().UTC(), Valid: true}
	(data)["enabled"] = false
	tx := s.DB.WithContext(ctx).Begin()
	model := new(models.User)
	err := tx.Model(&model).
		Where("id = ? AND deleted_by is null", int_id).
		Updates(data).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error in deleting model")
	}
	tx.Commit()
	return nil
}

func (s *UserService) GetUsers(ctx context.Context) (*[]dtos.UserResponse, error) {
	slice := make([]models.User, 0)
	s.DB.Where("enabled = ?", true).Find(&slice)
	res := make([]dtos.UserResponse, 0)
	for _, user := range slice {
		response, _:= TypeComverter[dtos.UserResponse](user)
		res = append(res, *response)
	}
	return &res, nil
}

func (s *UserService) GetProfile(ctx context.Context) (*models.User, error) {
	user := models.User{}
	id := ctx.Value("user_id")
	s.DB.Model(&user).Where("id = ?", id).First(&user)
	return &user, nil
}

func (s *UserService) GetFollowers(ctx context.Context) (*[]dtos.UserResponse, error) {
	user_id := ctx.Value("user_id")
	tx := s.DB.WithContext(ctx).Begin()
	user_followers := []models.UserFollowers{}
	err := tx.Model(&models.UserFollowers{}).Where("user_id = ? AND deleted_at is null", user_id).Find(&user_followers).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	follower := []dtos.UserResponse{}
	for _, user_follower := range user_followers {
		follower_id := user_follower.FollowerId
		follower_profile := models.User{}
		err = tx.Model(&models.User{}).Where("id = ? AND deleted_at is null", follower_id).First(&follower_profile).Error
		if err != nil {
			return nil, err
		}
		follower_res, _:= TypeComverter[dtos.UserResponse](follower_profile)
		follower = append(follower, *follower_res)
	}
	tx.Commit()
	return &follower, nil
}

func (s *UserService) GetFollowings(ctx context.Context) (*[]dtos.UserResponse, error) {
	user_id := ctx.Value("user_id")
	tx := s.DB.WithContext(ctx).Begin()
	user_followings := []models.UserFollowers{}
	err := tx.Model(&models.UserFollowers{}).Where("follower_id = ? AND deleted_at is null", user_id).Find(&user_followings).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	followings := []models.User{}
	for _, user_following := range user_followings {
		following_id := user_following.UserId
		user := models.User{}
		err = tx.Model(&models.User{}).Where("id = ? AND deleted_at is null", following_id).First(&user).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		followings = append(followings, user)
	}
	followings_res := []dtos.UserResponse{}
	for _, following := range followings {
		res, _:= TypeComverter[dtos.UserResponse](following)
		followings_res = append(followings_res, *res)
	}
	tx.Commit()
	return &followings_res, nil
}