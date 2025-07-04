package models

import (
	"time"

	"github.com/gocql/gocql"
)

type Image struct {
	ImageID               gocql.UUID `json:"image_id"`
	DayBucket             string     `json:"day_bucket,omitempty"`
	UploadedAt            time.Time  `json:"uploaded_at"`
	UserID                gocql.UUID `json:"user_id"`
	Username              string     `json:"username"`
	UserProfilePictureURL string     `json:"user_profile_picture_url"`
	ImageURL              string     `json:"image_url"`
	Title                 string     `json:"title"`
}
