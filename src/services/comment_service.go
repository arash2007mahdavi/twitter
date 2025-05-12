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
	comment.TweetId = tweet_id1
	err := tx.Model(&models.Comment{}).Create(&comment).Error
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error in creating comment")
	}
	tx.Commit()
	return TypeComverter[dtos.CommentResponse](comment)
}