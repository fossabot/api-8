package types

import "time"

type UserProfile struct {
	UserID                   int64     `json:"-"`
	ActivationToken          *string   `json:"-"`
	ActivationTokenExpiresAt time.Time `json:"-"`
	Name                     string    `json:"name"`
	Email                    string    `json:"email"`
	Address                  *string   `json:"address"`
	Phone                    *string   `json:"phone"`
}
