package redis

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

func (r *Redis) Del(key ...string) error {
	cmd := "DEL"
	args := make([]interface{}, len(key))
	for i, v := range key {
		args[i] = v
	}
	_, err := r.Exec(cmd, args...)
	return err
}

func (r *Redis) Exists(key ...string) (int, error) {
	cmd := "EXISTS"
	args := make([]interface{}, len(key))
	for i, v := range key {
		args[i] = v
	}
	return redis.Int(r.Exec(cmd, args...))
}

func (r *Redis) Expire(key string, expiration time.Duration) error {
	cmd := "EXPIRE"
	_, err := r.Exec(cmd, key, int64(expiration.Seconds()))
	return err
}

func (r *Redis) Expireat(key string, timestamp int64) error {
	cmd := "EXPIREAT"
	_, err := r.Exec(cmd, key, timestamp)
	return err
}

func (r *Redis) TTL(key string) (int64, error) {
	cmd := "TTL"
	ttl, err := r.Exec(cmd, key)
	if err != nil {
		return -1, err
	}
	return ttl.(int64), nil
}
