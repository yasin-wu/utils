package redis

import (
	"strconv"

	"github.com/gomodule/redigo/redis"
)

type Z struct {
	Score  int64
	Member string
}

func (this *Client) ZAdd(key string, z ...Z) error {
	cmd := "ZADD"
	args := make([]interface{}, len(z)*2+1)
	args[0] = key
	for i, v := range z {
		args[i*2+1] = v.Score
		args[i*2+2] = v.Member
	}
	_, err := this.Exec(cmd, args...)
	return err
}

func (this *Client) ZRange(key string, start, end int, withScores bool) ([]Z, error) {
	cmd := "ZRANGE"
	args := make([]interface{}, 3, 4)
	args[0] = key
	args[1] = start
	args[2] = end
	if withScores {
		args = append(args, "WITHSCORES")
	}
	var err error
	values, err := redis.Values(this.Exec(cmd, args...))
	if err != nil {
		return nil, err
	}
	return this.handleValues(values, withScores), nil
}

func (this *Client) ZRangeByScore(key, min, max string, withScores, limit bool, offset, count int) ([]Z, error) {
	cmd := "ZRANGEBYSCORE"
	args := make([]interface{}, 3, 7)
	args[0] = key
	args[1] = min
	args[2] = max
	if withScores {
		args = append(args, "WITHSCORES")
	}
	if limit {
		args = append(args, "LIMIT", offset, count)
	}
	var err error
	values, err := redis.Values(this.Exec(cmd, args...))
	if err != nil {
		return nil, err
	}
	return this.handleValues(values, withScores), nil
}

func (this *Client) ZRemrangEByScore(key, min, max string) error {
	cmd := "ZREMRANGEBYSCORE"
	args := make([]interface{}, 3)
	args[0] = key
	args[1] = min
	args[2] = max
	_, err := this.Exec(cmd, args...)
	return err
}

func (this *Client) handleValues(values []interface{}, withScores bool) []Z {
	var err error
	var redisZs []Z
	if withScores {
		for i := 0; i < len(values)/2; i++ {
			var z Z
			z.Member = this.readString(values[i*2])
			z.Score, err = strconv.ParseInt(this.readString(values[i*2+1]), 10, 64)
			if err != nil {
				continue
			}
			redisZs = append(redisZs, z)
		}
	} else {
		for i := 0; i < len(values); i++ {
			var z Z
			z.Member = this.readString(values[i])
			redisZs = append(redisZs, z)
		}
	}
	return redisZs
}
