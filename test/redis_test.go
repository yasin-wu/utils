package test

import (
	"testing"
	"time"

	js "github.com/bitly/go-simplejson"

	"github.com/yasin-wu/utils/redis"
)

var key = "123"

func TestRedis(t *testing.T) {
	conf := &redis.Config{
		Host:     "47.108.155.25:6379",
		PassWord: "yasinwu",
		DB:       2,
	}
	cli, err := redis.New(conf)
	if err != nil {
		t.Error(err)
		return
	}
	cli.DB = 1
	j := js.New()
	j.Set("a", 123)
	err = cli.Set(key, j, time.Minute)
	if err != nil {
		t.Error(err)
		return
	}
	err = cli.Del(key)
	if err != nil {
		t.Error(err)
		return
	}
}
