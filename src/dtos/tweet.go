package dtos

type TweetCreate struct {
	Title   string `json:"title" binding:"required,min=2,max=50"`
	Message string `json:"message" binding:"required,min=10,max=1000"`
}

type TweetUpdate struct {
	Title   string `json:"title,omitempty" binding:"omitempty,min=2,max=50"`
	Message string `json:"message,omitempty" binding:"omitempty,min=10,max=1000"`
}

type TweetResponse struct {
	Id       int               `json:"id"`
	UserId   int               `json:"user_id,omitempty"`
	User     *UserResponse     `json:"user,omitempty"`
	Title    string            `json:"title"`
	Message  string            `json:"message"`
	Comments []CommentResponse `json:"comments,omitempty"`
	Likes    []UserResponse    `json:"likes,omitempty"`
	Dislikes []UserResponse    `json:"dislikes,omitempty"`
	Files    []FileResponse    `json:"files,omitempty"`
}
