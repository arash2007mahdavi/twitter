package models

type User struct {
	BaseModel
	Username     string    `gorm:"size:30;not null;unique"`
	Firstname    string    `gorm:"size:20;null"`
	Lastname     string    `gorm:"size:40;null"`
	MobileNumber string    `json:"mobile_number" gorm:"size:11;not null;unique"`
	Password     string    `gorm:"size:50000;not null"`
	Enabled      bool      `gorm:"default:true"`
	Tweets       []Tweet   `gorm:"foreignKey:UserId"`
	Comments     []Comment `gorm:"foreignKey:UserId"`
	Followers    []User    `gorm:"many2many:follows;joinForeignKey:FollowingID;JoinReferences:FollowerID"`
	Followings   []User    `gorm:"many2many:follows;joinForeignKey:FollowerID;JoinReferences:FollowingID"`
}

type Tweet struct {
	BaseModel
	Title    string    `gorm:"size:50;not null"`
	Message  string    `gorm:"size:1000;not null"`
	UserId   int       `json:"user_id"`
	User     User      `gorm:"foreignKey:UserId"`
	Comments []Comment `gorm:"foreignKey:TweetId"`
}

type Comment struct {
	BaseModel
	TweetId int    `json:"tweet_id"`
	Tweet   Tweet  `gorm:"foreignKey:TweetId"`
	UserId  int    `json:"user_id"`
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
