package docker

import (
	"fmt"
	"time"

	"github.com/ory/dockertest"
)

var pool *dockertest.Pool

type checkFunc func() error

// Configure docker using available environment variables.
func Configure() {
	var err error
	pool, err = dockertest.NewPool("")
	if err != nil {
		panic(fmt.Sprintf("error while configuring docker pool, err: %s", err))
	}
}

func startContainer(image, tag string, env []string) (*dockertest.Resource, error) {
	resource, err := pool.Run(image, tag, env)
	if err != nil {
		return nil, err
	}
	return resource, nil
}

func waitContainer(f checkFunc) error {
	pool.MaxWait = 1 * time.Minute
	return pool.Retry(f)
}
