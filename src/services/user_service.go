package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"
	"twitter/src/database"
	"twitter/src/database/models"
	"twitter/src/dtos"
	"twitter/src/logger"
	"twitter/src/metrics"

	"gorm.io/gorm"
)

func TypeConverter[T any](req any) (*T, error) {
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
	data, err := TypeConverter[models.User](req)
	if err != nil {
		return nil, fmt.Errorf("error in comverting model")
	}
	tx := s.DB.WithContext(ctx).Begin()
	user_test := models.User{}
	err = tx.Model(&models.User{}).Where("username = ? AND enabled is true", req.Username).First(&user_test).Error
	if err == nil {
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(*data).String(), "Create", "Failed").Inc()
		return nil, fmt.Errorf("the username used by someone else")
	}
	user_test = models.User{}
	err = tx.Model(&models.User{}).Where("mobile_number = ? AND enabled is true", req.MobileNumber).First(&user_test).Error
	if err == nil {
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(*data).String(), "Create", "Failed").Inc()
		return nil, fmt.Errorf("the mobile number used by someone else")
	}
	err = tx.Create(data).Error
	if err != nil {
		tx.Rollback()
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(*data).String(), "Create", "Failed").Inc()
		return nil, fmt.Errorf("error in creating new model")
	}
	tx.Commit()
	metrics.DbCalls.WithLabelValues(reflect.TypeOf(*data).String(), "Create", "Success").Inc()
	return TypeConverter[dtos.UserResponse](data)
}

func (s *UserService) Update(ctx context.Context, req *dtos.UserUpdate) (*dtos.UserResponse, error) {
	data, err := TypeConverter[map[string]interface{}](req)
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
		Where("id = ? AND enabled is true", ctx.Value("user_id")).
		Updates(*data).Error
	if err != nil {
		tx.Rollback()
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(model).String(), "Update", "Failed").Inc()
		return nil, fmt.Errorf("error in updating model")
	}
	user := models.User{}
	tx.Model(&models.User{}).Where("id = ?", int_id).First(&user)
	tx.Commit()
	metrics.DbCalls.WithLabelValues(reflect.TypeOf(model).String(), "Update", "Success").Inc()
	return TypeConverter[dtos.UserResponse](user)
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
		Where("id = ? AND enabled is true", int_id).
		Updates(data).Error
	if err != nil {
		tx.Rollback()
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(model).String(), "Delete", "Failed").Inc()
		return fmt.Errorf("error in deleting model")
	}
	tx.Commit()
	metrics.DbCalls.WithLabelValues(reflect.TypeOf(model).String(), "Delete", "Success").Inc()
	return nil
}

func (s *UserService) GetUsers(ctx context.Context) (*[]dtos.UserResponse, error) {
	users := []dtos.UserResponse{}
	err := s.DB.Table("users").Where("enabled = ?", true).Scan(&users).Error
	if err != nil {
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.User{}).String(), "GetUsers", "Failed").Inc()
		return nil, err
	}
	metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.User{}).String(), "GetUsers", "Success").Inc()
	return &users, nil
}

func (s *UserService) GetProfile(ctx context.Context) (*models.User, error) {
	user := models.User{}
	id := ctx.Value("user_id")
	err := s.DB.Preload("Tweets", "enabled = ?", true).Preload("Tweets.Comments").
				Preload("Tweets.Likes").Preload("Tweets.Dislikes").
				Preload("Comments", "enabled = ?", true).Preload("Comments.Tweet").
				Preload("Comments.Likes").Preload("Comments.Dislikes").Preload("Followings").
				Preload("Followers").Preload("TweetLikes").Preload("TweetDislikes").
				Preload("CommentLikes").Preload("CommentDislikes").Model(&user).
				Where("id = ?", id).First(&user).Error
	if err != nil {
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(user).String(), "GetProfile", "Failed").Inc()
		return nil ,err
	}
	metrics.DbCalls.WithLabelValues(reflect.TypeOf(user).String(), "GetProfile", "Success").Inc()
	return &user, nil
}

func (s *UserService) GetFollowers(ctx context.Context) (*[]dtos.UserResponse, error) {
	user_id := ctx.Value("user_id")
	tx := s.DB.WithContext(ctx).Begin()

	user := models.User{}
	err := tx.Preload("Followers").Model(&models.User{}).
			  Where("id = ? AND enabled is true", user_id).First(&user).Error
	if err != nil {
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(user).String(), "GetFollowers", "Failed").Inc()
		return nil, err
	}

	followers := []dtos.UserResponse{}
	for _, follower := range user.Followers {
		res, _:= TypeConverter[dtos.UserResponse](follower)
		followers = append(followers, *res)
	}
	metrics.DbCalls.WithLabelValues(reflect.TypeOf(user).String(), "GetFollowers", "Success").Inc()
	return &followers, nil
}

func (s *UserService) GetFollowings(ctx context.Context) (*[]dtos.UserResponse, error) {
	user_id := ctx.Value("user_id")
	tx := s.DB.WithContext(ctx).Begin()

	user := models.User{}
	err := tx.Preload("Followings").Model(&models.User{}).
			  Where("id = ? AND enabled is true", user_id).First(&user).Error
	if err != nil {
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(user).String(), "GetFollowings", "Failed").Inc()
		return nil, err
	}

	followings := []dtos.UserResponse{}
	for _, following := range user.Followings {
		res, _:= TypeConverter[dtos.UserResponse](following)
		followings = append(followings, *res)
	}
	metrics.DbCalls.WithLabelValues(reflect.TypeOf(user).String(), "GetFollowings", "Success").Inc()
	return &followings, nil
}

func (s *UserService) Follow(ctx context.Context) error {
	user_id := ctx.Value("user_id")
	target_id := ctx.Value("target_id")
	tx := s.DB.WithContext(ctx).Begin()
	var user, target models.User
	err := tx.Model(&models.User{}).Where("id = ? AND enabled is true", user_id).First(&user).Error
	if err != nil {
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(user).String(), "Follow", "Failed").Inc()
		return err
	}
	err = tx.Model(&models.User{}).Where("id = ? AND enabled is true", target_id).First(&target).Error
	if err != nil {
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(user).String(), "Follow", "Failed").Inc()
		return err
	}
	err = tx.Model(&user).Association("Followings").Append(&target)
	if err != nil {
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(user).String(), "Follow", "Failed").Inc()
		return err
	}
	tx.Commit()
	metrics.DbCalls.WithLabelValues(reflect.TypeOf(user).String(), "Follow", "Success").Inc()
	return nil
}

func (s *UserService) UnFollow(ctx context.Context) error {
	user_id := ctx.Value("user_id")
	target_id := ctx.Value("target_id")
	tx := s.DB.WithContext(ctx).Begin()
	var user, target models.User
	err := tx.Model(&models.User{}).Where("id = ? AND enabled is true", user_id).First(&user).Error
	if err != nil {
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(user).String(), "UnFollow", "Failed").Inc()
		return err
	}
	err = tx.Model(&models.User{}).Where("id = ? AND enabled is true", target_id).First(&target).Error
	if err != nil {
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(user).String(), "UnFollow", "Failed").Inc()
		return err
	}
	err = tx.Model(&user).Association("Followings").Delete(&target)
	if err != nil {
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(user).String(), "UnFollow", "Failed").Inc()
		return err
	}
	tx.Commit()
	metrics.DbCalls.WithLabelValues(reflect.TypeOf(user).String(), "UnFollow", "Success").Inc()
	return nil
}