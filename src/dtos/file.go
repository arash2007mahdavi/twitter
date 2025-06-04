package dtos

import "mime/multipart"

type FileFormRequest struct {
	File *multipart.FileHeader `json:"file" form:"file" swaggerignore:"true"`
}

type UploadFileRequest struct {
	FileFormRequest
	Description string `json:"description,omitempty" form:"description"`
}

type CreateFileRequest struct {
	Name        string `json:"name"`
	Directory   string `json:"directory"`
	Description string `json:"description"`
	MimeType    string `json:"mime_type"`
	TweetId     *int   `json:"tweet_id,omitempty" gorm:"omitempty"`
	CommentId   *int   `json:"comment_id,omitempty" gorm:"omitempty"`
}

type UpdateFileRequest struct {
	Description string `json:"description,omitempty"`
}

type FileResponse struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Directory   string `json:"directory"`
	Description string `json:"description,omitempty"`
	MimeType    string `json:"mime_type"`
	Base64      string `json:"base64,omitempty"`
	TweetId     int    `json:"tweet_id,omitempty"`
	CommentId   int    `json:"comment_id,omitempty"`
}
