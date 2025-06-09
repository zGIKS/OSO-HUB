package models

import (
	"time"

	"github.com/gocql/gocql"
)

// Role constants for user roles
const (
	RoleUser   = "user"
	RoleAdmin  = "admin"
	RoleBanned = "banned"
)

type User struct {
	UserID            gocql.UUID `json:"user_id"`
	Username          string     `json:"username"`
	Email             string     `json:"email"`
	PasswordHash      string     `json:"password_hash"`
	ProfilePictureURL string     `json:"profile_picture_url"`
	Bio               string     `json:"bio"`
	Role              string     `json:"role"`
	CreatedAt         time.Time  `json:"created_at"`
}
