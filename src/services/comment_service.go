package services

import (
	"context"
	"fmt"
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
	tweet_id := ctx.Value("tweet_id")
	tweet_id1, _:= strconv.Atoi(tweet_id.(string))
	tx := s.Database.WithContext(ctx).Begin()
	comment, _:= TypeComverter[models.Comment](req)
	comment.CreatedBy = user_id.(int)
	comment.UserId = user_id.(int)
	user := dtos.UserResponse{}
	err := tx.Model(&models.User{}).Where("id = ?", comment.UserId).First(&user).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	comment.TweetId = tweet_id1
	tweet := dtos.TweetResponse{}
	err = tx.Model(&models.Tweet{}).Where("id = ?", comment.TweetId).First(&tweet).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = tx.Model(&models.Comment{}).Create(&comment).Error
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error in creating comment")
	}
	tx.Commit()
	comment_res, _:= TypeComverter[dtos.CommentResponse](comment)
	comment_res.User = user
	comment_res.Tweet = tweet
	return comment_res, nil
}