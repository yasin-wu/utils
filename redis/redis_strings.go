package redis

import (
	"encoding/json"
	"time"

	"github.com/gomodule/redigo/redis"
)

func (this *Client) Set(key string, value interface{}, expiration time.Duration) error {
	var err error
	cmd := "SET"
	args := make([]interface{}, 2, 5)
	args[0] = key
	args[1], err = json.Marshal(value)
	if err != nil {
		return err
	}
	if expiration > 0 {
		args = append(args, "EX", int64(expiration.Seconds()))
	}
	_, err = this.Exec(cmd, args...)
	return err
}

func (this *Client) Get(key string) ([]byte, error) {
	cmd := "GET"
	return redis.Bytes(this.Exec(cmd, key))
}
