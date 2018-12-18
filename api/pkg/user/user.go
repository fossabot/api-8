package user

import (
	"gitlab.com/pinterkode/pinterkode/api/pkg/database"
	"gitlab.com/pinterkode/pinterkode/api/pkg/model"
	"gitlab.com/pinterkode/pinterkode/api/pkg/utils/crypto"
)

func createUser(username, password, email string) (*model.User, error) {
	hashedPassword, err := crypto.CreateStringHash(password, 10)
	if err != nil {
		return nil, err
	}

	user := &model.User{
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

func saveUser(user *model.User) error {
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
