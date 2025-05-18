package dtos

import "mime/multipart"

type FileFormRequest struct {
	File *multipart.FileHeader `json:"file" binding:"required" swaggerignor:"true"`
}

type UploadFileRequest struct {
	FileFormRequest
	Description string `json:"description" binding:"required"`
}

type CreateFileRequest struct {
	Name        string `json:"name"`
	Directory   string `json:"directory"`
	Description string `json:"description"`
	MimeType    string `json:"mime_type"`
}

type UpdateFileRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type FileResponse struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Directory   string `json:"directory"`
	Description string `json:"description"`
	MimeType    string `json:"mime_type"`
}
