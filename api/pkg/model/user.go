package model

import (
	"time"

	"github.com/rs/xid"
)

type User struct {
	ID                       int64      `json:"id"`
	Username                 string     `json:"username"`
	EncryptedPassword        string     `json:"encrypted_password"`
	Email                    string     `json:"email"`
	Active                   bool       `json:"active"`
	ActivationToken          string     `json:"activation_token"`
	ActivationTokenExpiresAt time.Time  `json:"activation_token_expires_at"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                *time.Time `json:"updated_at"`
	DeletedAt                *time.Time `json:"deleted_at"`
}

func (u *User) GenerateActivationToken() {
	u.ActivationToken = xid.New().String()
}
