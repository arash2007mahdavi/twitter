package models

type User struct {
	BaseModel
	Username     string    `json:"username" gorm:"size:30;not null;unique"`
	Firstname    string    `json:"firstname" gorm:"size:20;null"`
	Lastname     string    `json:"lastname" gorm:"size:40;null"`
	MobileNumber string    `json:"mobile_number" gorm:"size:11;not null;unique"`
	Password     string    `json:"password" gorm:"size:50000;not null"`
	Enabled      bool      `json:"enabled" gorm:"default:true"`
	Tweets       []Tweet   `gorm:"foreignKey:UserId"`
	Comments     []Comment `gorm:"foreignKey:UserId"`
	Followers    []User `gorm:"many2many:follows;joinForeignKey:FollowingID;JoinReferences:FollowerID"`
    Following    []User `gorm:"many2many:follows;joinForeignKey:FollowerID;JoinReferences:FollowingID"`
}

// type UserFollowers struct {
// 	BaseModel
// 	UserId     int
// 	User       User `gorm:"foreignKey:UserId"`
// 	FollowerId int
// 	Follower   User `gorm:"foreignKey:FollowerId"`
// }

// type TweetLikes struct {
// 	BaseModel
// 	TweetId int
// 	Tweet   Tweet `gorm:"foreignKey:TweetId"`
// 	UserId  int
// 	User    User `gorm:"foreignKey:UserId"`
// }

type Tweet struct {
	BaseModel
	Title   string `gorm:"size:50;not null"`
	Message string `gorm:"size:1000;not null"`
	UserId  int    `json:"user_id" gorm:"not null"`
	User    User   `gorm:"foreignKey:UserId"`
	Comments []Comment `gorm:"foreignKey:TweetId"`
}

// type CommentLikes struct {
// 	BaseModel
// 	CommentId int
// 	Comment   Comment `gorm:"foreignKey:CommentId"`
// 	UserId    int
// 	User      User `gorm:"foreignKey:UserId"`
// }

type Comment struct {
	BaseModel
	TweetId int    `json:"tweet_id" gorm:"not null"`
	Tweet   Tweet  `gorm:"foreignKey:TweetId"`
	UserId  int    `json:"user_id" gorm:"not null"`
	User    User   `gorm:"foreignKey:UserId"`
	Message string `gorm:"size:1000;not null"`
}

type File struct {
	BaseModel
	Name        string `gorm:"size:100;type:string;not null"`
	Directory   string `gorm:"size:100;type:string;not null"`
	Description string `gorm:"size:500;type:string;not null"`
	MimeType    string `gorm:"size:20;type:string;not null"`
}
