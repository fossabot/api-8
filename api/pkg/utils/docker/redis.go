package docker

import (
	"fmt"

	"github.com/go-redis/redis"
)

// RunRedis create and run redis server.
// This method waits until the container is ready (container is up and can handle ping).
// It returns a dsn and function to destroy the container.
// This method is safe to be called multiple times.
// Each call spawns a new container with a new ip.
func RunRedis(version string) (string, func()) {
	if version == "" {
		version = "latest"
	}

	resource, err := startContainer("redis", version, nil)
	if err != nil {
		panic(fmt.Sprintf("fail to start redis container, err: %s", err))
	}

	var dsn = fmt.Sprintf("redis://%s:%s", resource.GetBoundIP("6379/tcp"), resource.GetPort("6379/tcp"))
	checkFunc := func() error {
		opts, err := redis.ParseURL(dsn)
		if err != nil {
			return err
		}
		client := redis.NewClient(opts)
		if val := client.Ping().Val(); val != "PONG" {
			return fmt.Errorf("fail to ping redis, got: %s", val)
		}
		return nil
	}
	err = waitContainer(checkFunc)
	if err != nil {
		panic(fmt.Sprintf("failed to wait redis to be ready, err: %s", err))
	}

	close := func() {
		resource.Close()
	}
	return dsn, close
}
