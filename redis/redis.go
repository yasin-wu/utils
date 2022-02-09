package redis

import (
	"errors"
	"reflect"
	"strings"
	"time"

	"log"

	"github.com/gomodule/redigo/redis"
)

const (
	defaultNetWork        = "tcp"
	defaultDB             = 0
	defaultMaxIdle        = 10
	defaultMaxActive      = 0
	defaultConnectTimeout = 5 * time.Second
	defaultReadTimeout    = 30 * time.Second
	defaultWriteTimeout   = 30 * time.Second
	defaultIdleTimeout    = 30 * time.Second
)

type Redis struct {
	DB int

	network        string
	password       string        //密码
	maxidle        int           //连接池最大空闲连接数
	maxactive      int           //连接池最大激活连接数
	connecttimeout time.Duration //连接超时
	readtimeout    time.Duration //读取超时
	writetimeout   time.Duration //写入超时
	pool           *redis.Pool
}

type Option func(client *Redis)

func New(host string, options ...Option) (*Redis, error) {
	if host == "" {
		return nil, errors.New("host is nil")
	}
	c := &Redis{}
	for _, f := range options {
		f(c)
	}
	checkConfig(c)
	redisDial := func() (redis.Conn, error) {
		conn, err := redis.Dial(
			strings.ToLower(c.network),
			host,
			redis.DialConnectTimeout(c.connecttimeout),
			redis.DialReadTimeout(c.readtimeout),
			redis.DialWriteTimeout(c.writetimeout),
		)
		if err != nil {
			log.Printf("连接redis失败:%s", err.Error())
			return nil, err
		}

		if c.password != "" {
			if _, err := conn.Do("AUTH", c.password); err != nil {
				conn.Close()
				log.Printf("redis认证失败:%s", err.Error())
				return nil, err
			}
		}

		_, err = conn.Do("SELECT", c.DB)
		if err != nil {
			conn.Close()
			log.Printf("redis选择数据库失败:%s", err.Error())
			return nil, err
		}

		return conn, nil
	}

	redisTestOnBorrow := func(conn redis.Conn, t time.Time) error {
		_, err := conn.Do("PING")
		if err != nil {
			log.Printf("从redis连接池取出的连接无效:%s", err.Error())
		}
		return err
	}

	pool := &redis.Pool{
		MaxIdle:      c.maxidle,
		MaxActive:    c.maxactive,
		IdleTimeout:  defaultIdleTimeout,
		Dial:         redisDial,
		TestOnBorrow: redisTestOnBorrow,
		Wait:         true,
	}
	c.pool = pool
	return c, nil
}

func checkConfig(c *Redis) {
	c.DB = defaultDB
	if c.network == "" {
		c.network = defaultNetWork
	}
	if c.maxidle == 0 {
		c.maxidle = defaultMaxIdle
	}
	if c.maxactive == 0 {
		c.maxactive = defaultMaxActive
	}
	if c.connecttimeout == 0 {
		c.connecttimeout = defaultConnectTimeout
	}
	if c.readtimeout == 0 {
		c.readtimeout = defaultReadTimeout
	}
	if c.writetimeout == 0 {
		c.writetimeout = defaultWriteTimeout
	}
}

func WithPassWord(passWord string) Option {
	return func(c *Redis) {
		c.password = passWord
	}
}

func WithNetWork(netWork string) Option {
	return func(c *Redis) {
		c.network = netWork
	}
}

func WithMaxIdle(maxIdle int) Option {
	return func(c *Redis) {
		c.maxidle = maxIdle
	}
}

func WithMaxActive(maxActive int) Option {
	return func(c *Redis) {
		c.maxactive = maxActive
	}
}

func WithConnectTimeout(connectTimeout time.Duration) Option {
	return func(c *Redis) {
		c.connecttimeout = connectTimeout
	}
}

func WithReadTimeout(readTimeout time.Duration) Option {
	return func(c *Redis) {
		c.readtimeout = readTimeout
	}
}

func WithWriteTimeout(writeTimeout time.Duration) Option {
	return func(c *Redis) {
		c.writetimeout = writeTimeout
	}
}

func (r *Redis) Exec(command string, args ...interface{}) (interface{}, error) {
	conn := r.pool.Get()
	defer conn.Close()
	_, err := conn.Do("SELECT", r.DB)
	if err != nil {
		return nil, err
	}
	return conn.Do(command, args...)
}

func (r *Redis) readString(value interface{}) string {
	var buffer []byte
	typeString := reflect.TypeOf(value).String()
	switch typeString {
	case "[]uint8":
		for _, v := range value.([]uint8) {
			buffer = append(buffer, v)
		}
	}
	return string(buffer)
}
