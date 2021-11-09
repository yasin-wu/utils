package test

import (
	"testing"

	"github.com/davecgh/go-spew/spew"

	"github.com/yasin-wu/utils/redis"
)

var key = "test-redis"

func TestRedis(t *testing.T) {
	conf := &redis.Config{
		Host:     "47.108.155.25:6379",
		PassWord: "yasinwu",
		DB:       0,
	}
	cli, err := redis.New(conf)
	if err != nil {
		t.Error(err)
		return
	}
	//cli.Set(key, "", time.Minute)
	spew.Dump(cli.TTL(key))
}
