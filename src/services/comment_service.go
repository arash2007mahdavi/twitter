package services

import (
	"context"
	"database/sql"
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

type CommentService struct {
	Database *gorm.DB
	Logger   logger.Logger
}

func NewCommentService() *CommentService {
	return &CommentService{
		Database: database.GetDB(),
		Logger:   logger.NewLogger(),
	}
}

func (s *CommentService) PostComment(ctx context.Context, req *dtos.CommentCreate) (*dtos.CommentResponse, error) {
	user_id := ctx.Value("user_id")
	user_id1, _:= user_id.(int)
	tweet_id := ctx.Value("tweet_id")
	tweet_id1, _:= strconv.Atoi(tweet_id.(string))
	tx := s.Database.WithContext(ctx).Begin()
	comment, _:= TypeConverter[models.Comment](req)
	comment.UserId = user_id1
	comment.CreatedBy = user_id1
	comment.TweetId = tweet_id1
	err := tx.Create(&comment).Error
	if err != nil {
		tx.Rollback()
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Comment{}).String(), "Create", "Failed").Inc()
		return nil, err
	}
	comment_1 := models.Comment{}
	err = tx.Preload("User").Preload("Tweet").Preload("Files").Preload("Tweet.User").Model(&models.Comment{}).
			 Where("id = ? AND enabled is true", comment.Id).First(&comment_1).Error
	if err != nil {
		tx.Rollback()
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Comment{}).String(), "Create", "Failed").Inc()
		return nil, err
	}
	comment_res, _:= TypeConverter[dtos.CommentResponse](comment_1)
	tx.Commit()
	metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Comment{}).String(), "Create", "Success").Inc()
	return comment_res, nil
}

func (s *CommentService) Update(ctx context.Context, req dtos.CommentUpdate) (dtos.CommentResponse, error) {
	comment_id := ctx.Value("comment_id")
	modified_by := ctx.Value("modified_by").(int)
	data, _:= TypeConverter[map[string]interface{}](req)
	(*data)["modified_by"] = sql.NullInt64{Int64: int64(modified_by), Valid: true}
	(*data)["modified_at"] = sql.NullTime{Time: time.Now().UTC(), Valid: true}
	tx := s.Database.WithContext(ctx).Begin()
	err := tx.Model(&models.Comment{}).Where("id = ? AND enabled is true", comment_id).Updates(&data).Error
	if err != nil {
		tx.Rollback()
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Comment{}).String(), "Update", "Failed").Inc()
		return dtos.CommentResponse{}, err
	}
	comment_res := models.Comment{}
	err = tx.Preload("Tweet").Preload("User").Preload("Likes").Preload("Files").
			 Preload("Dislikes").Model(&models.Comment{}).
			 Where("id = ? AND enabled is true", comment_id).First(&comment_res).Error
	if err != nil {
		tx.Rollback()
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Comment{}).String(), "Update", "Failed").Inc()
		return dtos.CommentResponse{}, err
	}
	res, _:= TypeConverter[dtos.CommentResponse](comment_res)
	tx.Commit()
	metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Comment{}).String(), "Update", "Success").Inc()
	return *res, nil
}

func (s *CommentService) Delete(ctx context.Context) error {
	comment_id := ctx.Value("comment_id")
	deleted_by := ctx.Value("deleted_by").(int)
	data := map[string]interface{}{}
	(data)["deleted_by"] = sql.NullInt64{Int64: int64(deleted_by), Valid: true}
	(data)["deleted_at"] = sql.NullTime{Time: time.Now().UTC(), Valid: true}
	(data)["enabled"] = false
	tx := s.Database.WithContext(ctx).Begin()
	err := tx.Model(&models.Comment{}).Where("id = ? AND enabled is true", comment_id).Updates(data).Error
	if err != nil {
		tx.Rollback()
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Comment{}).String(), "Delete", "Failed").Inc()
		return err
	}
	tx.Commit()
	metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Comment{}).String(), "Delete", "Success").Inc()
	return nil
}

func (s *CommentService) GetCommentById(ctx context.Context) (*models.Comment, error) {
	comment_id := ctx.Value("comment_id")
	tx := s.Database.WithContext(ctx).Begin()
	comment := models.Comment{}
	err := tx.Preload("Tweet", "enabled = ?", true).Preload("User", "enabled = ?", true).
			  Preload("Tweet.User", "enabled = ?", true).Preload("Tweet.Likes").Preload("Files").
			  Preload("Tweet.Dislikes").Preload("Likes").Preload("Dislikes").
			  Model(&models.Comment{}).Where("id = ? AND enabled is true", comment_id).First(&comment).Error
	if err != nil {
		tx.Rollback()
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Comment{}).String(), "GetComment", "Failed").Inc()
		return nil, err
	}
	tx.Commit()
	metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Comment{}).String(), "GetComment", "Success").Inc()
	return &comment, nil
}

func (s *CommentService) GetComments(ctx context.Context) ([]dtos.CommentResponse, error) {
	user_id := ctx.Value("user_id")
	tx := s.Database.WithContext(ctx).Begin()
	comments := []models.Comment{}
	err := tx.Preload("Tweet", "enabled = ?", true).Preload("User", "enabled = ?", true).
			  Preload("Tweet.User", "enabled = ?", true).Preload("Tweet.Likes").
			  Preload("Tweet.Dislikes").Preload("Likes").Preload("Dislikes").Preload("Files").
			  Model(&models.Comment{}).Where("user_id = ? AND enabled is true", user_id).Find(&comments).Error
	if err != nil {
		tx.Rollback()
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Comment{}).String(), "GetComments", "Failed").Inc()
		return []dtos.CommentResponse{}, err
	}
	tx.Commit()
	comments_res := []dtos.CommentResponse{}
	for _, comment := range comments {
		res, _:= TypeConverter[dtos.CommentResponse](comment)
		comments_res = append(comments_res, *res)
	}
	metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Comment{}).String(), "GetComments", "Success").Inc()
	return comments_res, nil
}

func (s *CommentService) LikeComment(ctx context.Context) error {
	user_id := ctx.Value("user_id")
	comment_id := ctx.Value("comment_id")
	tx := s.Database.WithContext(ctx).Begin()
	user := models.User{}
	err := tx.Model(&models.User{}).Where("id = ? AND enabled is true", user_id).First(&user).Error
	if err != nil {
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Comment{}).String(), "LikeComment", "Failed").Inc()
		return err
	}
	comment := models.Comment{}
	err = tx.Model(&models.Comment{}).Where("id = ? AND enabled is true", comment_id).First(&comment).Error
	if err != nil {
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Comment{}).String(), "LikeComment", "Failed").Inc()
		return err
	}
	err = tx.Model(&user).Association("CommentLikes").Append(&comment)
	if err != nil {
		tx.Rollback()
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Comment{}).String(), "LikeComment", "Failed").Inc()
		return err
	}
	tx.Commit()
	metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Comment{}).String(), "LikeComment", "Success").Inc()
	return nil
}

func (s *CommentService) DislikeComment(ctx context.Context) error {
	user_id := ctx.Value("user_id")
	comment_id := ctx.Value("comment_id")
	tx := s.Database.WithContext(ctx).Begin()
	user := models.User{}
	err := tx.Model(&models.User{}).Where("id = ? AND enabled is true", user_id).First(&user).Error
	if err != nil {
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Comment{}).String(), "DislikeComment", "Failed").Inc()
		return err
	}
	comment := models.Comment{}
	err = tx.Model(&models.Comment{}).Where("id = ? AND enabled is true", comment_id).First(&comment).Error
	if err != nil {
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Comment{}).String(), "DislikeComment", "Failed").Inc()
		return err
	}
	err = tx.Model(&user).Association("CommentLikes").Delete(&comment)
	if err != nil {
		tx.Rollback()
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Comment{}).String(), "DislikeComment", "Failed").Inc()
		return err
	}
	err = tx.Model(&user).Association("CommentDislikes").Append(&comment)
	if err != nil {
		tx.Rollback()
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Comment{}).String(), "DislikeComment", "Failed").Inc()
		return err
	}
	tx.Commit()
	metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Comment{}).String(), "DislikeComment", "Success").Inc()
	return nil
}