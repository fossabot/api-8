package testhelper

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/stretchr/testify/suite"
	"gitlab.com/pinterkode/pinterkode/api/pkg/database"
	"gitlab.com/pinterkode/pinterkode/api/pkg/utils/docker"
	"gitlab.com/pinterkode/pinterkode/api/pkg/utils/logger"
)

type Suite struct {
	suite.Suite
	suite.TearDownAllSuite

	NeedDB    bool
	destroyDB func()
}

func (s *Suite) SetupTest() {
	logger.SurpressLog()

	docker.Configure()
	if s.NeedDB {
		s.setupDB()
	}
}

func (s *Suite) TearDownSuite() {
	if s.NeedDB {
		database.Shutdown()
		s.destroyDB()
	}
}

func (s *Suite) RespBodyEqual(body []byte, other map[string]interface{}) {
	p, err := json.Marshal(other)
	s.Nil(err)
	s.Equal(string(body), string(p))
}

func (s *Suite) setupDB() {
	var dbURL string
	dbURL, s.destroyDB = docker.RunPostgres("11")
	database.ConfigureTest(dbURL)

	s.runDBMigration(dbURL)
}

func (s *Suite) runDBMigration(dbURL string) {
	wd, _ := os.Getwd()

	for n := 0; wd != "/" || n < 10; n++ {
		testDir := path.Join(wd, "database", "Rakefile")

		if _, err := os.Stat(testDir); !os.IsNotExist(err) {
			cmd := exec.Command("rake", "db:migrate")
			cmd.Env = append(cmd.Env, fmt.Sprintf("DB_URL=%s", dbURL))
			cmd.Dir = path.Join(wd, "database")
			_, err := cmd.CombinedOutput()
			if err != nil {
				s.FailNow("fail running db migration")
			}
			break
		}

		// go up one level
		wd = path.Join(wd, "..")
	}
}
