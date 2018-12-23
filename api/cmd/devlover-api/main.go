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
		logrus.WithError(err).Fatalln("failed to parse environment variables")
	}

	if err := database.Configure(&database.Config{
		Master: &database.DBConf{
			URL:          conf.DbURL,
			ConnLifetime: 60 * time.Minute,
			MaxIdleConns: 2,
			MaxOpenConns: 5,
		},
	}); err != nil {
		logrus.WithError(err).Fatalln("failed to configure database")
	}

	if err := server.Run(conf.ListenAddr, conf.Production); err != nil {
		logrus.WithField("msg", err.Error()).Warnln("server died")
	}
}
