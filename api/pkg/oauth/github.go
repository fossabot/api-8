package oauth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/devlover-id/api/pkg/config"
)

const (
	githubAccessTokenURL = "https://github.com/login/oauth/access_token"
	githubUserProfileURL = "https://api.github.com/user"
)

var githubClient = &http.Client{
	Timeout: 10 * time.Second,
}

func getGithubAccessToken(code string) (string, error) {
	reqBody := url.Values{}
	reqBody.Set("code", code)
	reqBody.Set("client_id", config.GithubClientID())
	reqBody.Set("client_secret", config.GithubClientSecret())

	req, err := http.NewRequest(http.MethodPost, githubAccessTokenURL, strings.NewReader(reqBody.Encode()))
	if err != nil {
		return "", err
	}

	resp, err := githubClient.Do(req)
	if err != nil {
		return "", err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	parsed, err := url.ParseQuery(string(respBody))
	if err != nil {
		return "", err
	}

	return parsed.Get("access_token"), nil
}

type GithubUserProfile struct {
	Username string  `json:"username"`
	Name     string  `json:"name"`
	Email    *string `json:"email"`
}

func (p *GithubUserProfile) UnmarshalJSON(b []byte) error {
	var temp struct {
		Username string  `json:"login"`
		Name     string  `json:"name"`
		Email    *string `json:"email"`
	}
	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}

	p.Username = temp.Username
	p.Name = temp.Name
	p.Email = temp.Email
	return nil
}

func getGithubUserProfile(accessToken string) (*GithubUserProfile, error) {
	req, err := http.NewRequest(http.MethodGet, githubUserProfileURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("authorization", fmt.Sprintf("bearer %s", accessToken))

	resp, err := githubClient.Do(req)
	if err != nil {
		return nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userProfile GithubUserProfile
	if err := json.Unmarshal(respBody, &userProfile); err != nil {
		return nil, err
	}

	return &userProfile, nil
}
