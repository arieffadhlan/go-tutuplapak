package dto

import "time"

type File struct {
	ID           string    `db:"id" `
	Url          string    `db:"url"`
	ThumbnailUrl string    `db:"thumbnail_url"`
	CreateAt     time.Time `db:"created_at"`
}

type FileResponse struct {
	ID           string `json:"fileId"`
	Url          string `json:"fileUri"`
	ThumbnailUrl string `json:"fileThumbnailUri"`
}
