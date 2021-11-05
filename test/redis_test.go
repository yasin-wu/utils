package test

import (
	"testing"

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
	/*var zs []redis.Z
	zs = append(zs, redis.Z{Score: 1, Member: "a"}, redis.Z{Score: 2, Member: "b"})
	err = cli.ZAdd(key, zs...)
	if err != nil {
		t.Error(err)
		return
	}*/
	err = cli.ZRemrangEByScore(key, "1", "(2")
	if err != nil {
		t.Error(err)
		return
	}
}
