package models

type User struct {
	BaseModel
	Username     string  `gorm:"size:30;not null;unique"`
	Firstname    string  `gorm:"size:20;null"`
	Lastname     string  `gorm:"size:40;null"`
	MobileNumber string  `json:"mobile_number" gorm:"size:11;not null;unique"`
	Password     string  `gorm:"size:50;not null"`
	Followers    []*User `gorm:"many2many:user_followers;joinForeignKey:UserId;joinReferences:FollowerId"`  
	Following    []*User `gorm:"many2many:user_followers;joinForeignKey:FollowerId;joinReferences:UserId"`
	Tweets       []Tweet
	Enabled      bool    `gorm:"default:true"`
}

type Tweet struct {
	BaseModel
	Title    string    `gorm:"size:50;not null"`
	Message  string    `gorm:"size:1000;not null"`
	UserId   int       `gorm:"not null"`
	User     *User      `gorm:"foreignKey:UserId"`
	Comments []Comment `gorm:"foreignKey:TweetId"`
}

type Comment struct {
	BaseModel
	TweetId int
	Tweet   *Tweet  `gorm:"foreignKey:TweetID"`
	Message string `gorm:"size:1000;not null"`
}
