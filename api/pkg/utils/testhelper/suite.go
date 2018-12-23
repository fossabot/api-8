package testhelper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/devlover-id/api/pkg/database"
	"github.com/devlover-id/api/pkg/utils/docker"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	suite.SetupAllSuite
	suite.TearDownAllSuite

	NeedDB    bool
	destroyDB func()
}

func (s *Suite) SetupSuite() {
	logrus.SetOutput(ioutil.Discard)

	docker.Configure()
	if s.NeedDB {
		s.setupDB()
	}
}

func (s *Suite) TearDownSuite() {
	if s.NeedDB && os.Getenv("DB_URL") == "" {
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
	if dbURL := os.Getenv("DB_URL"); dbURL == "" {
		dbURL, s.destroyDB = docker.RunPostgres("11")
		database.ConfigureTest(dbURL)
		s.runDBMigration(dbURL)
	} else {
		database.ConfigureTest(dbURL)
	}
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
