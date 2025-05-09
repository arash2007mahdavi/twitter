package models

type User struct {
    BaseModel
    Username     string  `gorm:"size:30;not null;unique"`
    Firstname    string  `gorm:"size:20;null"`
    Lastname     string  `gorm:"size:40;null"`
    MobileNumber string  `json:"mobile_number" gorm:"size:11;not null;unique"`
    Password     string  `gorm:"size:50;not null"`
    Followers    []int   `gorm:"type:integer[];default:'{}'"`
    Following    []int   `gorm:"type:integer[];default:'{}'"`
    Tweets       []Tweet `gorm:"foreignKey:UserID"` // ارتباط کاربر به توییت‌ها از طریق فیلد UserID
    Enabled      bool    `gorm:"default:true"`
}

type Tweet struct {
    BaseModel
    Title    string    `gorm:"size:50;not null"`
    Message  string    `gorm:"size:1000;not null"`
    UserID   uint      `gorm:"not null"`
    User     User      
    Comments []Comment `gorm:"foreignKey:TweetID"`
}

type Comment struct {
    BaseModel
    TweetID int
    Tweet   Tweet   `gorm:"foreignKey:TweetID"`
    Message string  `gorm:"size:1000;not null"`
}