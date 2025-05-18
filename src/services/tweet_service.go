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
	user := dtos.UserResponse{}
	err = tx.Model(&models.User{}).Where("id = ?", id_creator).First(&user).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = tx.Model(&models.Tweet{}).Create(&tweet).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	tweet_res, _:= TypeComverter[dtos.TweetResponse](tweet)
	tweet_res.User = user
	return tweet_res, nil
}

func (s *TweetService) GetTweetByID(ctx context.Context) (*dtos.TweetResponse, error) {
	tweet_id, _:= strconv.Atoi(ctx.Value("tweet_id").(string))
	tx := s.Database.WithContext(ctx).Begin()
	var tweet dtos.TweetResponse
	err := tx.Model(&models.Tweet{}).Where("id = ? AND deleted_by is null", tweet_id).First(&tweet).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	user_id := tweet.UserId
	user := dtos.UserResponse{}
	err = tx.Model(&models.User{}).Where("id = ?", user_id).First(&user).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tweet.User = user
	comments := []dtos.CommentResponse{}
	tx.Model(&models.Comment{}).Where("tweet_id = ? AND deleted_at is null", tweet_id).Find(&comments)
	tweet.Comments = comments
	tx.Commit()
	return &tweet, nil
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
		tweets_comments := []models.Comment{}
		err = tx.Model(&models.Comment{}).Where("tweet_id = ? AND deleted_at is null", tweet.Id).Find(&tweets_comments).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		comments_response := []dtos.CommentResponse{}
		for _, comment := range tweets_comments {
			res, _:= TypeComverter[dtos.CommentResponse](comment)
			comments_response = append(comments_response, *res)
		}
		res.Comments = comments_response
		tweets_response = append(tweets_response, *res)
	}
	tx.Commit()
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

func (s *TweetService) GetFollowingsTweets(ctx context.Context) ([]dtos.TweetResponse, error) {
	user_id := ctx.Value("user_id")
	tx := s.Database.WithContext(ctx).Begin()
	followings := []models.UserFollowers{}
	err := tx.Model(&models.UserFollowers{}).Where("follower_id = ? AND deleted_at is null", user_id).Find(&followings).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	followings_tweets := []dtos.TweetResponse{}
	for _, userFollower := range followings {
		tweets := []models.Tweet{}
		err = tx.Model(&models.Tweet{}).Where("user_id = ? AND deleted_at is null", userFollower.UserId).Find(&tweets).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		tweets_res := []dtos.TweetResponse{}
		for _, tweet:= range tweets {
			tweet_res, _:= TypeComverter[dtos.TweetResponse](tweet)
			tweets_res = append(tweets_res, *tweet_res)
		}
		for _, tweet := range tweets_res {
			comments := []models.Comment{}
			err = tx.Model(&models.Comment{}).Where("tweet_id = ? AND deleted_at is null", tweet.Id).Find(&comments).Error
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			comments_res := []dtos.CommentResponse{}
			for _, comment := range comments {
				comment_res, _:= TypeComverter[dtos.CommentResponse](comment)
				comments_res = append(comments_res, *comment_res)
			}
			tweet.Comments = comments_res
			followings_tweets = append(followings_tweets, tweet)
		}
	} 
	tx.Commit()
	return followings_tweets, nil
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