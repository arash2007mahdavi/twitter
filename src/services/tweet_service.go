package services

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
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
	err = tx.Create(&tweet).Error
	if err != nil {
		return nil, err
	}
	tweet_2 := models.Tweet{}
	err = tx.Preload("User").Model(&models.Tweet{}).Where("id = ? AND deleted_at is null", tweet.Id).First(&tweet_2).Error
	if err != nil {
		return nil, err
	}
	tweet_res := dtos.TweetResponse{}
	res, _:= TypeComverter[dtos.TweetResponse](tweet_2)
	tweet_res = *res
	tx.Commit()
	return &tweet_res, nil
}

func (s *TweetService) GetTweetByID(ctx context.Context) (*dtos.TweetResponse, error) {
	tweet_id, _:= strconv.Atoi(ctx.Value("tweet_id").(string))
	tx := s.Database.WithContext(ctx).Begin()
	var tweet models.Tweet
	err := tx.Preload("User").Preload("Comments").Preload("Comments.User").Model(&models.Tweet{}).Where("id = ? AND deleted_by is null", tweet_id).First(&tweet).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tweet_res, _:= TypeComverter[dtos.TweetResponse](tweet)
	tx.Commit()
	return tweet_res, nil
}

func (s *TweetService) GetTweets(ctx context.Context) ([]dtos.TweetResponse, error) {
	user_id := ctx.Value("user_id").(int)
	tx := s.Database.WithContext(ctx).Begin()
	tweets := []models.Tweet{}
	err := tx.Preload("User").Preload("Comments").Preload("Comments.User").Model(&models.Tweet{}).Where("user_id = ? AND deleted_by is null", user_id).Find(&tweets).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	tweets_res := []dtos.TweetResponse{}
	for _, tweet := range tweets {
		res, _:= TypeComverter[dtos.TweetResponse](tweet)
		tweets_res = append(tweets_res, *res)
	}
	return tweets_res, nil
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
	tweet := dtos.TweetResponse{}
	err = tx.Preload("User").Preload("Comments").Preload("Comments.User").Model(&models.Tweet{}).Where("id = ? AND deleted_by is null", tweet_id).First(&tweet).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return &tweet, nil
} 

func (s *TweetService) Delete(ctx context.Context) error {
	tweet_id := ctx.Value("tweet_id")
	deleted_by := ctx.Value("deleted_by").(int)
	data := map[string]interface{}{}
	(data)["deleted_by"] = sql.NullInt64{Int64: int64(deleted_by), Valid: true}
	(data)["deleted_at"] = sql.NullTime{Time: time.Now().UTC(), Valid: true}
	tx := s.Database.WithContext(ctx).Begin()
	err := tx.Model(&models.Tweet{}).Where("id = ? AND deleted_at is null", tweet_id).Updates(data).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("tweet couldnt delete")
	}
	data_comments := map[string]interface{}{}
	(data_comments)["deleted_by"] = sql.NullInt64{Int64: int64(deleted_by), Valid: true}
	(data_comments)["deleted_at"] = sql.NullTime{Time: time.Now().UTC(), Valid: true}
	err = tx.Model(&models.Comment{}).Where("tweet_id = ? AND deleted_at is null", tweet_id).Updates(data_comments).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("comments couldnt delete")
	}
	tx.Commit()
	return nil
}

func (s *TweetService) GetFollowingsTweets(ctx context.Context) ([]dtos.TweetResponse, error) {
	user_id := ctx.Value("user_id")
	tx := s.Database.WithContext(ctx).Begin()
	user := models.User{}
	err := tx.Preload("Followings").Preload("Followings.Tweets").Model(&models.User{}).Where("id = ? AND deleted_at is null", user_id).First(&user).Error
	if err != nil {
		return nil, err
	}
	followings := user.Followings
	followings_tweets := []models.Tweet{}
	for _, following := range followings {
		followings_tweets = append(followings_tweets, following.Tweets...)
	}
	tx.Commit()
	res := []dtos.TweetResponse{}
	for _, tweet := range followings_tweets {
		res_tweet, _:= TypeComverter[dtos.TweetResponse](tweet)
		res = append(res, *res_tweet)
	}
	return res, nil
}

func selectRandomTweets(slice []models.Tweet, count int) []models.Tweet {
    if count > len(slice) {
        count = len(slice)
    }

    rand.Shuffle(len(slice), func(i, j int) { slice[i], slice[j] = slice[j], slice[i] })
    return slice[:count]
}
func (s *TweetService) TweetExplore(ctx context.Context) (*[]dtos.TweetResponse, error) {
	tx := s.Database.WithContext(ctx).Begin()
	tweets := []models.Tweet{}
	err := tx.Model(&models.Tweet{}).Where("deleted_at is null").Find(&tweets).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	five_tweets := selectRandomTweets(tweets, 5)
	tweets_res := []dtos.TweetResponse{}
	for _, tweet := range five_tweets {
		res, _:= TypeComverter[dtos.TweetResponse](tweet)
		comments := []models.Comment{}
		err = tx.Model(&models.Comment{}).Where("tweet_id = ? AND deleted_at is null", res.Id).Find(&comments).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		comments_res := []dtos.CommentResponse{}
		for _, comment := range comments {
			res, _:= TypeComverter[dtos.CommentResponse](comment)
			comments_res = append(comments_res, *res)
		}
		res.Comments = comments_res
		tweets_res = append(tweets_res, *res)
	}
	tx.Commit()
	return &tweets_res, nil
}