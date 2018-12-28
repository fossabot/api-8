package testhelper

import (
	"github.com/devlover-id/api/pkg/database"
	"github.com/devlover-id/api/pkg/types"
	"github.com/devlover-id/api/pkg/utils/crypto"
	"github.com/icrowley/fake"
)

func (s *Suite) NewUser() (*types.User, string) {
	tx, err := database.NewTransaction()
	s.Nil(err)
	s.NotNil(tx)

	plainPassword := fake.Password(8, 12, true, true, true)
	hashedPassword, err := crypto.CreateStringHash(plainPassword, 1)
	s.Nil(err)

	var user types.User
	err = tx.Query(&user, `
		insert into users (
			username,
			encrypted_password,
			authentication_token
		)
		values (?, ?, ?)
		returning *
	`, fake.UserName(), hashedPassword, fake.CharactersN(10))
	s.Nil(err)

	var profile types.UserProfile
	err = tx.Query(&profile, `
		insert into user_profile(
			user_id,
			activation_token,
			name,
			email
		)
		values (?, ?, ?, ?)
		returning *
	`, user.ID, fake.CharactersN(10), fake.FullName(), fake.EmailAddress())
	s.Nil(err)

	err = tx.Commit()
	s.Nil(err)

	user.Profile = profile
	return &user, plainPassword
}
