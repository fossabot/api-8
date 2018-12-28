package testhelper

import (
	"io/ioutil"
	"os"

	"github.com/devlover-id/api/pkg/database"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	suite.SetupAllSuite
	suite.TearDownAllSuite

	NeedDB bool
}

func (s *Suite) SetupSuite() {
	logrus.SetOutput(ioutil.Discard)

	if s.NeedDB {
		s.setupDB()
	}
}

func (s *Suite) TearDownSuite() {
	if s.NeedDB {
		database.Shutdown()
	}
}

func (s *Suite) setupDB() {
	dbURL := os.Getenv("DB_URL")
	db, err := database.ConfigureTest(dbURL)
	s.Nil(err)
	s.Nil(db.Ping())
}
