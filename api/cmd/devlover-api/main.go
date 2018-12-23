package main

import (
	"time"

	"github.com/caarlos0/env"
	"github.com/devlover-id/api/pkg/database"
	"github.com/devlover-id/api/pkg/server"
	"github.com/sirupsen/logrus"
)

func main() {
	var conf config
	if err := env.Parse(&conf); err != nil {
		logrus.Fatalln("failed to parse environment variables", err)
	}

	if err := database.Configure(&database.Config{
		Master: &database.DBConf{
			URL:          conf.ListenAddr,
			ConnLifetime: 60 * time.Minute,
			MaxIdleConns: 2,
			MaxOpenConns: 5,
		},
	}); err != nil {
		logrus.Fatalln("failed to configure database", err)
	}

	if err := server.Run(conf.ListenAddr, false); err != nil {
		logrus.WithField("msg", err.Error()).Println("server exitted")
	}
}
