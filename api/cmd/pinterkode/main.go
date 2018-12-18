package main

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/pinterkode/pinterkode/api/pkg/database"
	"gitlab.com/pinterkode/pinterkode/api/pkg/server"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if err := database.Configure(&database.Config{
		Master: &database.DBConf{
			URL:          os.Getenv("DB_URL"),
			ConnLifetime: 60 * time.Minute,
			MaxIdleConns: 5,
			MaxOpenConns: 5,
		},
	}); err != nil {
		fmt.Println(err)
	}

	addr := "localhost:8080"
	if err := server.Run(addr, false); err != nil {
		fmt.Println(err)
	}
}
