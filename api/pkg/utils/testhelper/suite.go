package testhelper

import (
	"encoding/json"
	"os"

	"github.com/stretchr/testify/suite"
	"gitlab.com/pinterkode/pinterkode/api/pkg/database"
	"gitlab.com/pinterkode/pinterkode/api/pkg/utils/logger"
)

type Suite struct {
	suite.Suite
	suite.TearDownAllSuite

	NeedDB bool
}

func (s *Suite) SetupTest() {
	logger.SurpressLog()
	if s.NeedDB {
		if os.Getenv("TEST_DB_URL") == "" {
			s.T().Fatal("TEST_DB_URL is empty")
		}
		database.ConfigureTest(os.Getenv("TEST_DB_URL"))
	}
}

func (s *Suite) TearDownSuite() {
	if s.NeedDB {
		database.Shutdown()
	}
}

func (s *Suite) RespBodyEqual(body []byte, other map[string]interface{}) {
	p, err := json.Marshal(other)
	s.Nil(err)
	s.Equal(string(body), string(p))
}
