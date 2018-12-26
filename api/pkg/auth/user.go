package auth

import (
	"github.com/devlover-id/api/pkg/database"
	"github.com/devlover-id/api/pkg/types"
	"github.com/devlover-id/api/pkg/utils/crypto"
	"github.com/rs/xid"
)

func createUser(username, password, name, email string) (*types.User, error) {
	hashedPassword, err := crypto.CreateStringHash(password, 10)
	if err != nil {
		return nil, err
	}
	activationToken := createActivationToken()
	authenticationToken := createAuthenticationToken()

	tx, err := database.NewTransaction()
	if err != nil {
		return nil, err
	}

	var newUser types.User
	createUserQuery := `
		insert into users (
			username,
			encrypted_password,
			authentication_token
		)
		values (?, ?, ?)
		returning *
	`
	if err := tx.Query(&newUser, createUserQuery, username, hashedPassword, authenticationToken); err != nil {
		return nil, err
	}

	var newProfile types.UserProfile
	createUserProfileQuery := `
		insert into user_profile(
			user_id,
			activation_token,
			name,
			email
		)
		values (?, ?, ?, ?)
		returning *
	`
	if err := tx.Query(&newProfile, createUserProfileQuery, newUser.ID, activationToken, name, email); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	newUser.Profile = newProfile
	return &newUser, nil
}

func createActivationToken() string {
	return xid.New().String()
}

func createAuthenticationToken() string {
	return xid.New().String()
}
