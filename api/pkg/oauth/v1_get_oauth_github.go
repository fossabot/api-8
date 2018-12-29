package oauth

import (
	"net/http"

	"github.com/devlover-id/api/pkg/api"
)

func V1GetGithub(req api.Request) api.Response {
	code := req.Raw().URL.Query().Get("code")
	if code == "" {
		return api.CodeOnlyResp(http.StatusBadRequest)
	}
	accessToken, err := getGithubAccessToken(code)
	if err != nil {
		return api.InternalServerErrResp(err)
	}
	userProfile, err := getGithubUserProfile(accessToken)
	if err != nil {
		return api.InternalServerErrResp(err)
	}
	registered, err := checkUsernameRegistered(userProfile.Username)
	if err != nil {
		return api.InternalServerErrResp(err)
	}
	return api.JSONResponse(http.StatusOK, map[string]interface{}{
		"github_profile": userProfile,
		"registered":     registered,
	})
}
