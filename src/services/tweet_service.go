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
	id_creator := ctx.Value("user_id").(int)
	tx := s.Database.WithContext(ctx).Begin()
	tweet, err := TypeComverter[models.Tweet](req)
	if err != nil {
		return nil, err
	}
	tweet.UserId = id_creator
	tweet.CreatedBy = id_creator
	err = tx.Model(&models.Tweet{}).Create(&tweet).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	user_tweet := models.UserTweet{
		UserId: id_creator, TweetId: tweet.Id,
	}
	err = tx.Model(&models.UserTweet{}).Create(&user_tweet).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return TypeComverter[dtos.TweetResponse](tweet)
}

func (s *TweetService) GetTweetByID(ctx context.Context) (*dtos.TweetResponse, error) {
	tweet_id, _:= strconv.Atoi(ctx.Value("tweet_id").(string))
	var tweet models.Tweet
	err := s.Database.WithContext(ctx).Model(&models.Tweet{}).Where("id = ? AND deleted_by is null", tweet_id).First(&tweet).Error
	if err != nil {
		return nil, err
	}
	return TypeComverter[dtos.TweetResponse](tweet)
}

func (s *TweetService) GetTweets(ctx context.Context) ([]dtos.TweetResponse, error) {
	user_id := ctx.Value("user_id").(int)
	tweets := []models.Tweet{}
	tx := s.Database.WithContext(ctx).Begin()
	err := tx.Model(&models.Tweet{}).Where("user_id = ? AND deleted_by is null", user_id).Find(&tweets).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tweets_response := []dtos.TweetResponse{}
	for _, tweet := range tweets {
		res, _:= TypeComverter[dtos.TweetResponse](tweet)
		tweets_response = append(tweets_response, *res)
	}
	return tweets_response, nil
}