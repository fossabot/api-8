package main

import (
	"time"

	"github.com/devlover-id/api/pkg/config"
	"github.com/devlover-id/api/pkg/database"
	"github.com/devlover-id/api/pkg/server"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := config.ParseEnv(); err != nil {
		logrus.WithError(err).Fatalln("failed to parse environment variables")
	}

	if err := database.Configure(&database.Config{
		Master: &database.DBConf{
			URL:          config.DbURL(),
			ConnLifetime: 60 * time.Minute,
			MaxIdleConns: 2,
			MaxOpenConns: 5,
		},
	}); err != nil {
		logrus.WithError(err).Fatalln("failed to configure database")
	}

	logrus.WithField("listen_addr", config.ListenAddr()).WithField("production", config.Production()).Infoln("starting api server")
	if err := server.Run(config.ListenAddr(), config.Production()); err != nil {
		logrus.WithField("msg", err.Error()).Warnln("server died")
	}
}
