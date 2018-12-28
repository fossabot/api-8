package auth

import (
	"errors"

	"github.com/devlover-id/api/pkg/database"
	"github.com/devlover-id/api/pkg/types"
	"github.com/devlover-id/api/pkg/utils/crypto"
	"github.com/rs/xid"
)

var errNilUserPassword = errors.New("nil user password")

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
		tx.Rollback()
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
		tx.Rollback()
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

// login authentication token, and nil error
// if given username and password are correct.
func login(username, password string) (string, error) {
	var user types.User
	var query = `
		select *
		from users
		where username = ?
	`
	if err := database.WriterQuery(&user, query, username); err != nil {
		return "", err
	}

	if user.EncryptedPassword == nil {
		return "", errNilUserPassword
	}
	if !crypto.ValidateHash(password, *user.EncryptedPassword) {
		return "", nil
	}

	return user.AuthenticationToken, nil
}
