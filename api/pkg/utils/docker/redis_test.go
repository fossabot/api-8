// +build integration

package docker

import (
	"testing"
	"time"

	"github.com/go-redis/redis"
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
	if err := client.Get("key").Err(); err != redis.Nil {
		t.Error("should be nil error")
	}

	// set something
	if err := client.Set("key", "something", time.Hour).Err(); err != nil {
		t.Error(err)
	}

	// read from key
	if val, err := client.Get("key").Result(); err != nil {
		t.Error(err)
	} else {
		if val != "something" {
			t.Error("val should equals something")
		}
	}

	// destroy container
	close()

	// read again
	if err := client.Get("key").Err(); err == nil {
		t.Error("should be error")
	}
}
