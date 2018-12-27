package testhelper

import (
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

func (s *Suite) setupDB() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL, s.destroyDB = docker.RunPostgres("11")
		s.runDbMigration(dbURL)
	}

	_, err := database.ConfigureTest(dbURL)
	s.Nil(err)
}

func (s *Suite) runDbMigration(dbURL string) {
	wd, _ := os.Getwd()
	for n := 0; wd != "/" || n < 8; n++ {
		testDir := path.Join(wd, "database", "Rakefile")

		if _, err := os.Stat(testDir); !os.IsNotExist(err) {
			cmds := []string{"db:drop", "db:create", "db:migrate"}
			for _, cmd := range cmds {
				c := exec.Command("rake", cmd)
				c.Env = append(c.Env, fmt.Sprintf("DB_URL=%s", dbURL))
				c.Dir = path.Join(wd, "database")
				out, err := c.CombinedOutput()
				if err != nil {
					s.T().Log("\n" + string(out))
					s.FailNow("failed running db migration")
				}
			}
			break
		}

		// go up one level
		wd = path.Join(wd, "..")
	}
}
