package services

import (
	"context"
	"strconv"
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
	err = tx.Preload("User").Preload("Tweet").Preload("Tweet.User").Model(&models.Comment{}).Where("id = ? AND deleted_at is null", comment.Id).First(&comment_1).Error
	if err != nil {
		return nil, err
	}
	comment_res, _:= TypeComverter[dtos.CommentResponse](comment_1)
	tx.Commit()
	return comment_res, nil
}