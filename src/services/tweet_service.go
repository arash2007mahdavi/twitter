package services

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
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
	tweet, err := TypeConverter[models.Tweet](req)
	if err != nil {
		return nil, err
	}
	tweet.UserId = id_creator
	tweet.CreatedBy = id_creator

	err = tx.Create(&tweet).Error
	if err != nil {
		tx.Rollback()
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Tweet{}).String(), "PostTweet", "Failed").Inc()
		return nil, err
	}
	tweet_2 := models.Tweet{}
	err = tx.Preload("User").Preload("Files").Model(&models.Tweet{}).Where("id = ? AND enabled is true", tweet.Id).First(&tweet_2).Error
	if err != nil {
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Tweet{}).String(), "PostTweet", "Failed").Inc()
		return nil, err
	}
	tweet_res := dtos.TweetResponse{}
	res, _:= TypeConverter[dtos.TweetResponse](tweet_2)
	tweet_res = *res
	tx.Commit()
	metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Tweet{}).String(), "PostTweet", "Success").Inc()
	return &tweet_res, nil
}

func (s *TweetService) GetTweetByID(ctx context.Context) (*models.Tweet, error) {
	tweet_id, _:= strconv.Atoi(ctx.Value("tweet_id").(string))
	tx := s.Database.WithContext(ctx).Begin()
	var tweet models.Tweet
	err := tx.Preload("User", "enabled", true).Preload("Comments", "enabled", true).Preload("Files").
			  Preload("Comments.User").Preload("Comments.Likes").Preload("Comments.Dislikes").
			  Preload("Likes").Preload("Dislikes").Model(&models.Tweet{}).
			  Where("id = ? AND enabled is true", tweet_id).First(&tweet).Error
	if err != nil {
		tx.Rollback()
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Tweet{}).String(), "GetTweet", "Failed").Inc()
		return nil, err
	}
	tx.Commit()
	metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Tweet{}).String(), "GetTweet", "Success").Inc()
	return &tweet, nil
}

func (s *TweetService) GetTweets(ctx context.Context) ([]dtos.TweetResponse, error) {
	user_id := ctx.Value("user_id").(int)
	tx := s.Database.WithContext(ctx).Begin()
	tweets := []models.Tweet{}
	err := tx.Preload("User", "enabled", true).Preload("Comments", "enabled", true).Preload("Files").
			  Preload("Comments.User", "enabled", true).Preload("Files").Preload("Likes").
			  Preload("Dislikes").Model(&models.Tweet{}).Where("user_id = ? AND enabled is true", user_id).Find(&tweets).Error
	if err != nil {
		tx.Rollback()
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Tweet{}).String(), "GetTweets", "Failed").Inc()
		return nil, err
	}
	tx.Commit()
	tweets_res := []dtos.TweetResponse{}
	for _, tweet := range tweets {
		res, _:= TypeConverter[dtos.TweetResponse](tweet)
		tweets_res = append(tweets_res, *res)
	}
	metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Tweet{}).String(), "GetTweets", "Success").Inc()
	return tweets_res, nil
}

func (s *TweetService) Update(ctx context.Context, req *dtos.TweetUpdate) (*dtos.TweetResponse, error) {
	data, _:= TypeConverter[map[string]interface{}](req)
	modified_by := ctx.Value("modified_by").(int)
	tweet_id := ctx.Value("tweet_id")
	(*data)["modified_by"] = sql.NullInt64{Int64: int64(modified_by), Valid: true}
	(*data)["modified_at"] = sql.NullTime{Time: time.Now().UTC(), Valid: true}
	tx := s.Database.WithContext(ctx).Begin()
	err := tx.Model(&models.Tweet{}).Where("id = ? AND enabled is true", tweet_id).Updates(*data).Error
	if err != nil {
		tx.Rollback()
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Tweet{}).String(), "Update", "Failed").Inc()
		return nil, err
	}
	tweet := dtos.TweetResponse{}
	err = tx.Preload("User").Preload("Comments", "enabled = ?", true).Preload("Files").
			 Preload("Comments.User").Preload("Likes").Preload("Dislikes").
			 Model(&models.Tweet{}).Where("id = ? AND enabled is true", tweet_id).First(&tweet).Error
	if err != nil {
		tx.Rollback()
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Tweet{}).String(), "Update", "Failed").Inc()
		return nil, err
	}
	tx.Commit()
	metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Tweet{}).String(), "Update", "Success").Inc()
	return &tweet, nil
} 

func (s *TweetService) Delete(ctx context.Context) error {
	tweet_id := ctx.Value("tweet_id")
	deleted_by := ctx.Value("deleted_by").(int)
	data := map[string]interface{}{}
	(data)["deleted_by"] = sql.NullInt64{Int64: int64(deleted_by), Valid: true}
	(data)["deleted_at"] = sql.NullTime{Time: time.Now().UTC(), Valid: true}
	(data)["enabled"] = false
	tx := s.Database.WithContext(ctx).Begin()
	err := tx.Model(&models.Tweet{}).Where("id = ? AND enabled is true", tweet_id).Updates(data).Error
	if err != nil {
		tx.Rollback()
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Tweet{}).String(), "Delete", "Failed").Inc()
		return fmt.Errorf("tweet couldnt delete")
	}
	data_comments := map[string]interface{}{}
	(data_comments)["deleted_by"] = sql.NullInt64{Int64: int64(deleted_by), Valid: true}
	(data_comments)["deleted_at"] = sql.NullTime{Time: time.Now().UTC(), Valid: true}
	(data_comments)["enabled"] = false
	err = tx.Model(&models.Comment{}).Where("tweet_id = ? AND enabled is true", tweet_id).Updates(data_comments).Error
	if err != nil {
		tx.Rollback()
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Tweet{}).String(), "Delete", "Failed").Inc()
		return fmt.Errorf("comments couldnt delete")
	}
	tx.Commit()
	metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Tweet{}).String(), "Delete", "Success").Inc()
	return nil
}

func (s *TweetService) GetFollowingsTweets(ctx context.Context) ([]dtos.TweetResponse, error) {
	user_id := ctx.Value("user_id")
	tx := s.Database.WithContext(ctx).Begin()
	user := models.User{}
	err := tx.Preload("Followings").Preload("Followings.Tweets", "enabled = ?", true).Preload("Followings.Files").
			  Preload("Followings.Tweets.Comments", "enabled = ?", true).Preload("Followings.Tweets.User", "enabled = ?", true).
			  Preload("Followings.Tweets.Likes").Preload("Followings.Tweets.Dislikes").Model(&models.User{}).
			  Where("id = ? AND enabled is true", user_id).First(&user).Error
	if err != nil {
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Tweet{}).String(), "GetFollowingsTweet", "Failed").Inc()
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
		res_tweet, _:= TypeConverter[dtos.TweetResponse](tweet)
		res = append(res, *res_tweet)
	}
	metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Tweet{}).String(), "GetFollowingsTweet", "Success").Inc()
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
	err := tx.Model(&models.Tweet{}).Where("enabled is true").Find(&tweets).Error
	if err != nil {
		tx.Rollback()
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Tweet{}).String(), "TweetExplore", "Failed").Inc()
		return nil, err
	}
	five_tweets := selectRandomTweets(tweets, 5)
	tweets_res := []dtos.TweetResponse{}
	for _, tweet := range five_tweets {
		res, _:= TypeConverter[dtos.TweetResponse](tweet)
		comments := []models.Comment{}
		err = tx.Model(&models.Comment{}).Where("tweet_id = ? AND enabled is true", res.Id).Find(&comments).Error
		if err != nil {
			tx.Rollback()
			metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Tweet{}).String(), "TweetExplore", "Failed").Inc()
			return nil, err
		}
		comments_res := []dtos.CommentResponse{}
		for _, comment := range comments {
			res, _:= TypeConverter[dtos.CommentResponse](comment)
			comments_res = append(comments_res, *res)
		}
		res.Comments = comments_res
		tweets_res = append(tweets_res, *res)
	}
	tx.Commit()
	metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Tweet{}).String(), "TweetExplore", "Success").Inc()
	return &tweets_res, nil
}

func (s *TweetService) LikeTweet(ctx context.Context) error {
	user_id := ctx.Value("user_id")
	tweet_id := ctx.Value("tweet_id")
	tx := s.Database.WithContext(ctx).Begin()
	user := models.User{}
	err := tx.Model(&models.User{}).Where("id = ? AND enabled is true", user_id).First(&user).Error
	if err != nil {
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Tweet{}).String(), "LikeTweet", "Failed").Inc()
		return err
	}
	tweet := models.Tweet{}
	err = tx.Model(&models.Tweet{}).Where("id = ? AND enabled is true", tweet_id).First(&tweet).Error
	if err != nil {
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Tweet{}).String(), "LikeTweet", "Failed").Inc()
		return err
	}
	err = tx.Model(&user).Association("TweetLikes").Append(&tweet)
	if err != nil {
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Tweet{}).String(), "LikeTweet", "Failed").Inc()
		return err
	}
	tx.Commit()
	metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Tweet{}).String(), "LikeTweet", "Success").Inc()
	return nil
}

func (s *TweetService) DislikeTweet(ctx context.Context) error {
	user_id := ctx.Value("user_id")
	tweet_id := ctx.Value("tweet_id")
	tx := s.Database.WithContext(ctx).Begin()
	user := models.User{}
	err := tx.Model(&models.User{}).Where("id = ? AND enabled is true", user_id).First(&user).Error
	if err != nil {
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Tweet{}).String(), "DislikeTweet", "Failed").Inc()
		return err
	}
	tweet := models.Tweet{}
	err = tx.Model(&models.Tweet{}).Where("id = ? AND enabled is true", tweet_id).First(&tweet).Error
	if err != nil {
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Tweet{}).String(), "DislikeTweet", "Failed").Inc()
		return err
	}
	err = tx.Model(&user).Association("TweetLikes").Delete(&tweet)
	if err != nil {
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Tweet{}).String(), "DislikeTweet", "Failed").Inc()
		tx.Rollback()
		return err
	}
	err = tx.Model(&user).Association("TweetDislikes").Append(&tweet)
	if err != nil {
		metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Tweet{}).String(), "DislikeTweet", "Failed").Inc()
		tx.Rollback()
		return err
	}
	tx.Commit()
	metrics.DbCalls.WithLabelValues(reflect.TypeOf(models.Tweet{}).String(), "DislikeTweet", "Success").Inc()
	return nil
}