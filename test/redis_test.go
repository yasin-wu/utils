package test

import (
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/yasin-wu/utils/redis"
)

var key = "test-redis"

func TestRedis(t *testing.T) {
	host := "47.108.155.25:6379"
	cli, err := redis.New(host, redis.WithPassWord("yasinwu"))
	if err != nil {
		t.Error(err)
		return
	}
	cli.Set(key, "", time.Minute)
	spew.Dump(cli.TTL(key))
}
