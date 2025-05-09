package services

import (
	"context"
	"twitter/src/database"
	"twitter/src/database/models"
	"twitter/src/dtos"
	"twitter/src/logger"

	"gorm.io/gorm"
)

type TweetService struct {
	Database *gorm.DB
	Logger logger.Logger
}

func NewTweetService() *TweetService {
	return &TweetService{
		Database: database.GetDB(),
		Logger: logger.NewLogger(),
	}
}

func (s *TweetService) NewTweet(ctx context.Context, req *dtos.TweetCreate) (*dtos.TweetResponse, error) {
	username := ctx.Value("username")
	tx := s.Database.WithContext(ctx).Begin()
	tweet, err := database.TypeComverter[models.Tweet](req)
	if err != nil {
		return nil, err
	}
	user := models.User{}
	err = tx.Model(&models.User{}).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	tweet.CreatedBy = user.Id
	tweet.UserID = uint(user.Id)
	err = tx.Model(&models.User{}).Where("username = ?", username).Association("Tweets").Append(tweet)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = tx.Model(&models.Tweet{}).Create(tweet).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return database.TypeComverter[dtos.TweetResponse](tweet)
}