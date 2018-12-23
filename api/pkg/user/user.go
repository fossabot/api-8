package user

import (
	"github.com/devlover-id/api/pkg/database"
	"github.com/devlover-id/api/pkg/types"
	"github.com/devlover-id/api/pkg/utils/crypto"
)

func createUser(username, password, email string) (*types.User, error) {
	hashedPassword, err := crypto.CreateStringHash(password, 10)
	if err != nil {
		return nil, err
	}

	user := &types.User{
		Username:          username,
		EncryptedPassword: hashedPassword,
		Email:             email,
	}
	user.GenerateActivationToken()

	if err := saveUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

func saveUser(user *types.User) error {
	query := `
		insert into users (
			username,
			encrypted_password,
			email,
			activation_token
		)
		values (?, ?, ?, ?)
		returning *
	`
	return database.WriterQuery(user, query, user.Username, user.EncryptedPassword, user.Email, user.ActivationToken)
}
