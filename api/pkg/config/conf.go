package config

import (
	"github.com/caarlos0/env"
)

type config struct {
	ListenAddr         string `env:"LISTEN_ADDR"`
	DbURL              string `env:"DB_URL"`
	Production         bool   `env:"PRODUCTION"`
	GithubClientID     string `env:"GITHUB_CLIENT_ID"`
	GithubClientSecret string `env:"GITHUB_CLIENT_SECRET"`
}

var conf config

func ParseEnv() error {
	return env.Parse(&conf)
}

func ListenAddr() string {
	return conf.ListenAddr
}

func DbURL() string {
	return conf.DbURL
}

func Production() bool {
	return conf.Production
}

func GithubClientID() string {
	return conf.GithubClientID
}

func GithubClientSecret() string {
	return conf.GithubClientSecret
}
