package test

import (
	"testing"

	js "github.com/bitly/go-simplejson"

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
	j := js.New()
	j.Set("a", 1)
	err = cli.Append(key, j)
	if err != nil {
		t.Error(err)
		return
	}
}
