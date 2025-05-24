package services

import (
	"context"
	"database/sql"
	"strconv"
	"time"
	"twitter/src/database"
	"twitter/src/database/models"
	"twitter/src/dtos"
	"twitter/src/logger"

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
	comment, _:= TypeComverter[models.Comment](req)
	comment.UserId = user_id1
	comment.CreatedBy = user_id1
	comment.TweetId = tweet_id1
	err := tx.Create(&comment).Error
	if err != nil {
		return nil, err
	}
	comment_1 := models.Comment{}
	err = tx.Preload("User").Preload("Tweet").Preload("Tweet.User").Model(&models.Comment{}).Where("id = ? AND enabled is true", comment.Id).First(&comment_1).Error
	if err != nil {
		return nil, err
	}
	comment_res, _:= TypeComverter[dtos.CommentResponse](comment_1)
	tx.Commit()
	return comment_res, nil
}

func (s *CommentService) Update(ctx context.Context, req dtos.CommentUpdate) (dtos.CommentResponse, error) {
	comment_id := ctx.Value("comment_id")
	modified_by := ctx.Value("modified_by").(int)
	data, _:= TypeComverter[map[string]interface{}](req)
	(*data)["modified_by"] = sql.NullInt64{Int64: int64(modified_by), Valid: true}
	(*data)["modified_at"] = sql.NullTime{Time: time.Now().UTC(), Valid: true}
	tx := s.Database.WithContext(ctx).Begin()
	err := tx.Model(&models.Comment{}).Where("id = ? AND enabled is true", comment_id).Updates(&data).Error
	if err != nil {
		tx.Rollback()
		return dtos.CommentResponse{}, err
	}
	comment_res := models.Comment{}
	err = tx.Preload("Tweet").Preload("User").Model(&models.Comment{}).Where("id = ? AND enabled is true", comment_id).First(&comment_res).Error
	if err != nil {
		tx.Rollback()
		return dtos.CommentResponse{}, err
	}
	res, _:= TypeComverter[dtos.CommentResponse](comment_res)
	tx.Commit()
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
		return err
	}
	tx.Commit()
	return nil
}