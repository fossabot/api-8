package types

import (
	"time"
)

type User struct {
	ID                  int64      `json:"id"`
	Username            string     `json:"username"`
	EncryptedPassword   *string    `json:"-"`
	GithubUsername      *string    `json:"github_username"`
	AuthenticationToken string     `json:"-"`
	Active              bool       `json:"active"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           *time.Time `json:"updated_at"`
	DeletedAt           *time.Time `json:"-"`

	Profile UserProfile `json:"profile"`
}
