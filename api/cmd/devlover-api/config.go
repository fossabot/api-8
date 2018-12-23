package main

// use upper case for env name
type config struct {
	ListenAddr string `env:"LISTEN_ADDR"`
	DbURL      string `env:"DB_URL"`
	Production bool   `env:"PRODUCTION"`
}
