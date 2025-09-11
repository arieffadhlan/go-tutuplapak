package dto

import "time"

type File struct {
	ID           string    `db:"id"`
	Url          string    `db:"url"`
	ThumbnailUrl string    `db:"thumbnail_url"`
	CreateAt     time.Time `db:"created_at"`
}
