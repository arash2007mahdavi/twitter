package models

type User struct {
	BaseModel
	Username        string    `json:"username,omitempty" gorm:"size:30;not null;unique"`
	Firstname       string    `json:"firstname,omitempty" gorm:"size:20;null"`
	Lastname        string    `json:"lastname,omitempty" gorm:"size:40;null"`
	MobileNumber    string    `json:"mobile_number,omitempty" gorm:"size:11;not null;unique"`
	Password        string    `json:"password,omitempty" gorm:"size:50000;not null"`
	Enabled         bool      `json:"enabled,omitempty" gorm:"default:true"`
	Tweets          []Tweet   `json:"tweets,omitempty" gorm:"foreignKey:UserId"`
	Comments        []Comment `json:"comments,omitempty" gorm:"foreignKey:UserId"`
	TweetLikes      []Tweet   `json:"tweet_likes,omitempty" gorm:"many2many:tweet_likes;joinForeignKey:UserId;JoinRefrences:TweetId"`
	TweetDislikes   []Tweet   `json:"tweet_dislikes,omitempty" gorm:"many2many:tweet_dislikes;joinForeignKey:UserId;JoinRefrences:TweetId"`
	CommentLikes    []Comment `json:"comment_likes,omitempty" gorm:"many2many:comment_likes;joinForeignKey:UserId;JoinRefrences:CommentId"`
	CommentDislikes []Comment `json:"comment_dislikes,omitempty" gorm:"many2many:comment_dislikes;joinForeignKey:UserId;JoinRefrences:CommentId"`
	Followers       []User    `json:"followers,omitempty" gorm:"many2many:follows;joinForeignKey:FollowingID;JoinReferences:FollowerID"`
	Followings      []User    `json:"followings,omitempty" gorm:"many2many:follows;joinForeignKey:FollowerID;JoinReferences:FollowingID"`
}

type Tweet struct {
	BaseModel
	Title    string    `json:"title,omitempty" gorm:"size:50;not null"`
	Message  string    `json:"message,omitempty" gorm:"size:1000;not null"`
	UserId   int       `json:"user_id,omitempty"`
	User     *User     `json:"user,omitempty" gorm:"foreignKey:UserId"`
	Comments []Comment `json:"comments,omitempty" gorm:"foreignKey:TweetId"`
	Files    []File    `json:"files,omitempty" gorm:"foreignKey:TweetId"`
	Likes    []User    `json:"likes,omitempty" gorm:"many2many:tweet_likes;joinForeignKey:TweetId;JoinRefrences:UserId"`
	Dislikes []User    `json:"dislikes,omitempty" gorm:"many2many:tweet_dislikes;joinForeignKey:TweetId;JoinRefrences:UserId"`
	Enabled  bool      `json:"enabled,omitempty" gorm:"default:true"`
}

type Comment struct {
	BaseModel
	TweetId  int    `json:"tweet_id,omitempty"`
	Tweet    *Tweet `json:"tweet,omitempty" gorm:"foreignKey:TweetId"`
	UserId   int    `json:"user_id,omitempty"`
	User     *User  `json:"user,omitempty" gorm:"foreignKey:UserId"`
	Message  string `json:"message,omitempty" gorm:"size:1000;not null"`
	Files    []File `json:"files,omitempty" gorm:"foreignKey:CommentId"`
	Likes    []User `json:"likes,omitempty" gorm:"many2many:comment_likes;joinForeignKey:CommentId;JoinRefrences:UserId"`
	Dislikes []User `json:"dislikes,omitempty" gorm:"many2many:comment_dislikes;joinForeignKey:CommentId;JoinRefrences:UserId"`
	Enabled  bool   `json:"enabled,omitempty" gorm:"default:true"`
}

type File struct {
	BaseModel
	Name        string   `gorm:"size:100;type:string;not null"`
	Directory   string   `gorm:"size:100;type:string;not null"`
	Description string   `gorm:"size:500;type:string;not null"`
	MimeType    string   `json:"mime_type" gorm:"size:20;type:string;not null"`
	Tweet       *Tweet   `json:"tweet,omitempty" gorm:"foreignKey:TweetId"`
	TweetId     *int      `json:"tweet_id,omitempty" gorm:"omitempty;null"`
	Comment     *Comment `json:"comment,omitempty" gorm:"foreignKey:CommentId"`
	CommentId   *int      `json:"comment_id,omitempty" gorm:"omitempty;null"`
}
