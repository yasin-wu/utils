package test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/yasin-wu/utils/redis"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func TestRedis(t *testing.T) {
	host := "47.108.155.25:6379"
	password := "yasinwu"
	key := "test-redis"
	cli, err := redis.New(host, redis.WithPassWord(password))
	if err != nil {
		log.Fatal(err)
	}
	_ = cli.Set(key, "", time.Minute)
	fmt.Println(cli.TTL(key))
}
