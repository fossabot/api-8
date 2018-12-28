package auth

import (
	"encoding/json"
	"net/http"

	"github.com/devlover-id/api/pkg/api"
	"github.com/devlover-id/api/pkg/types"
	"github.com/devlover-id/api/pkg/utils/testhelper"
	"github.com/icrowley/fake"
)

func (s *AuthTestSuite) TestPostRegister() {
	password := fake.Password(8, 12, true, true, true)
	payload := &v1PostRegisterPayload{
		Name:            fake.FullName(),
		Username:        fake.UserName(),
		Password:        password,
		PasswordConfirm: password,
		Email:           fake.EmailAddress(),
	}
	req := api.NewDummyRequest().
		SetJSONBody(payload)

	resp := V1PostRegister(testhelper.NewContext(), req)
	s.Equal(resp.StatusCode(), http.StatusCreated)

	user := types.User{}
	user.Profile = types.UserProfile{}

	err := json.Unmarshal(resp.Body(), &user)
	s.Nil(err)
	s.NotZero(user.ID)
	s.NotZero(user.Username)
	s.Zero(user.EncryptedPassword)
	s.Zero(user.GithubUsername)
	s.Zero(user.AuthenticationToken)
	s.False(user.Active)
	s.NotZero(user.CreatedAt)
	s.Nil(user.UpdatedAt)
	s.Nil(user.DeletedAt)
	s.Zero(user.Profile.UserID)
	s.Zero(user.Profile.ActivationToken)
	s.Zero(user.Profile.ActivationTokenExpiresAt)
	s.NotZero(user.Profile.Name)
	s.NotZero(user.Profile.Email)
	s.Nil(user.Profile.Address)
	s.Nil(user.Profile.Phone)
}
