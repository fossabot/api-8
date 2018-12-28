package auth

import (
	"net/http"

	"github.com/devlover-id/api/pkg/api"
)

func (s *AuthTestSuite) TestV1PostLogin_WithValidationError() {
	resp := V1PostLogin(api.NewDummyRequest())
	s.Equal(resp.StatusCode(), http.StatusBadRequest)
	s.RespBodyEqual(resp.Body(), map[string]string{
		"password": "password tidak boleh kosong",
		"username": "username tidak boleh kosong",
	})
}

func (s *AuthTestSuite) TestV1PostLogin_WithCorrectCredentials() {
	user, pass := s.NewUser()
	req := api.NewDummyRequest().SetJSONBody(map[string]string{
		"username": user.Username,
		"password": pass,
	})

	resp := V1PostLogin(req)
	s.Equal(resp.StatusCode(), http.StatusOK)
	s.RespBodyEqual(resp.Body(), map[string]string{
		"token": user.AuthenticationToken,
	})
}

func (s *AuthTestSuite) TestV1PostLogin_WithIncorrectCredentials() {
	user, _ := s.NewUser()
	req := api.NewDummyRequest().SetJSONBody(map[string]string{
		"username": user.Username,
		"password": "wrong password",
	})
	resp := V1PostLogin(req)
	s.Equal(resp.StatusCode(), http.StatusForbidden)
}
