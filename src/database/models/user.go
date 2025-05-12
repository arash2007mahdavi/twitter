package models

import "twitter/src/dtos"

type User struct {
	BaseModel
	Username     string `gorm:"size:30;not null;unique"`
	Firstname    string `gorm:"size:20;null"`
	Lastname     string `gorm:"size:40;null"`
	MobileNumber string `json:"mobile_number" gorm:"size:11;not null;unique"`
	Password     string `gorm:"size:50000;not null"`
	Enabled      bool   `gorm:"default:true"`
}

type UserFollowers struct {
	BaseModel
	UserId     int
	User       User `gorm:"foreignKey:UserId"`
	FollowerId int
	Follower   User `gorm:"foreignKey:FollowerId"`
}

type Tweet struct {
	BaseModel
	Title    string   `gorm:"size:50;not null"`
	Message  string   `gorm:"size:1000;not null"`
	UserId   int      `json:"user_id" gorm:"not null"`
	User     User     `gorm:"foreignKey:UserId"`
	Comments []dtos.CommentResponse `gorm:"foreignKey:TweetId"`
}

type Comment struct {
	BaseModel
	TweetId int    `json:"tweet_id" gorm:"not null"`
	Tweet   Tweet  `gorm:"foreignKey:TweetId"`
	UserId  int    `json:"user_id" gorm:"not null"`
	User    User   `gorm:"foreignKey:UserId"`
	Message string `gorm:"size:1000;not null"`
}