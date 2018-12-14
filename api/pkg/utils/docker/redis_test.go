// +build integration

package docker

import (
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

func TestRunRedis(t *testing.T) {
	Configure()

	dsn, close := RunRedis("")
	opts, err := redis.ParseURL(dsn)
	if err != nil {
		t.Fatal(err)
	}

	client := redis.NewClient(opts)

	// try get empty key
	err = client.Get("key").Err()
	assert.Equal(t, err, redis.Nil)

	// set something
	err = client.Set("key", "something", time.Hour).Err()
	assert.Nil(t, err)

	// read from key
	val, err := client.Get("key").Result()
	assert.Nil(t, err)
	assert.Equal(t, val, "something")

	// destroy container
	close()

	// read again
	err = client.Get("key").Err()
	assert.NotNil(t, err)
}
