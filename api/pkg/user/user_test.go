package user

import (
	"github.com/icrowley/fake"
	"gitlab.com/pinterkode/pinterkode/api/pkg/model"
)

func (s *UserTestSuite) TestSaveUser() {
	user := &model.User{
		Username:          fake.UserName(),
		EncryptedPassword: fake.CharactersN(10),
		Email:             fake.EmailAddress(),
		ActivationToken:   fake.CharactersN(10),
	}
	err := saveUser(user)
	s.Nil(err)

	s.True(user.ID > 0)
	s.False(user.Active)
	s.NotZero(user.ActivationTokenExpiresAt)
	s.NotZero(user.CreatedAt)
	s.Nil(user.UpdatedAt)
	s.Nil(user.DeletedAt)
}

func (s *UserTestSuite) TestCreateUser() {
	user, err := createUser(fake.UserName(), fake.SimplePassword(), fake.EmailAddress())
	s.Nil(err)
	s.NotNil(user)
	s.NotZero(user.ActivationToken)
	s.NotZero(user.ActivationTokenExpiresAt)
}
