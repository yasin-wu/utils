package test

import (
	"testing"

	mredis "github.com/gomodule/redigo/redis"

	"github.com/yasin-wu/utils/redis"
)

func TestRedis(t *testing.T) {
	conf := &redis.Config{
		Host:     "192.168.131.135:6379",
		PassWord: "1qazxsw21201",
	}
	cli, err := redis.New(conf)
	if err != nil {
		t.Error(err)
		return
	}
	res, err := cli.Exec("SET", "123", "456")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)
	resp, err := mredis.Bytes(cli.Exec("GET", "123"))
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(string(resp))
}
