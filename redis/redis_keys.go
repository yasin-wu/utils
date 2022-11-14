package redis

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

func (this *Client) Del(key ...string) error {
	cmd := "DEL"
	args := make([]interface{}, len(key))
	for i, v := range key {
		args[i] = v
	}
	_, err := this.Exec(cmd, args...)
	return err
}

func (this *Client) Exists(key ...string) (int, error) {
	cmd := "EXISTS"
	args := make([]interface{}, len(key))
	for i, v := range key {
		args[i] = v
	}
	return redis.Int(this.Exec(cmd, args...))
}

func (this *Client) Expire(key string, expiration time.Duration) error {
	cmd := "EXPIRE"
	_, err := this.Exec(cmd, key, int64(expiration.Seconds()))
	return err
}
