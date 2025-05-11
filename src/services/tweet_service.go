package services

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"
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

func (s *TweetService) Update(ctx context.Context, req *dtos.TweetUpdate) (*dtos.TweetResponse, error) {
	data, _:= TypeComverter[map[string]interface{}](req)
	modified_by := ctx.Value("modified_by").(int)
	tweet_id := ctx.Value("tweet_id")
	(*data)["modified_by"] = sql.NullInt64{Int64: int64(modified_by), Valid: true}
	(*data)["modified_at"] = sql.NullTime{Time: time.Now().UTC(), Valid: true}
	tx := s.Database.WithContext(ctx).Begin()
	err := tx.Model(&models.Tweet{}).Where("id = ? AND deleted_by is null", tweet_id).Updates(*data).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tweet := models.Tweet{}
	tx.Model(&models.Tweet{}).Where("id = ? AND deleted_by is null", tweet_id).First(&tweet)
	tx.Commit()
	return TypeComverter[dtos.TweetResponse](tweet)
} 

func (s *TweetService) Delete(ctx context.Context) error {
	tweet_id := ctx.Value("tweet_id")
	deleted_by := ctx.Value("deleted_by").(int)
	data := map[string]interface{}{}
	(data)["deleted_by"] = sql.NullInt64{Int64: int64(deleted_by), Valid: true}
	(data)["deleted_at"] = sql.NullTime{Time: time.Now().UTC(), Valid: true}
	tx := s.Database.WithContext(ctx).Begin()
	err := tx.Model(&models.Tweet{}).Where("id = ? AND deleted_by is null", tweet_id).Updates(data).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("user doesnt exists")
	}
	tx.Commit()
	return nil
}