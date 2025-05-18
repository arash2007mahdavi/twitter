package dtos

type CommentCreate struct {
	Message string `json:"message" binding:"required,min=5,max=300"`
}

type CommentUpdate struct {
	Message string `json:"message,omitempty" binding:"omitempty,min=5,max=300"`
}

type CommentResponse struct {
	TweetId int           `json:"tweet_id"`
	Tweet   TweetResponse `json:"tweet"`
	UserId  int           `json:"user_id"`
	User    UserResponse  `json:"user"`
	Message string        `json:"message"`
}
