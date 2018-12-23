package main

type config struct {
	ListenAddr string `env:"listen_addr"`
	DbURL      string `env:"db_url"`
}
