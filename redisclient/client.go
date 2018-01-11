package redisclient

import (
	"errors"

	"log"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/sudiyi/sdy/utils"
)

const (
	Success int = 1 // 成功

	DefaultRedisDb int = 0

	DefaultMaxIdle     int           = 3                 // 空闲连接的最大数目
	DefaultMaxActive   int           = 1000              // 给定时间内最大连接数，为0则连接数没有限制
	DefaultMaxWaitTime time.Duration = 180 * time.Second // Redis最大等待时间
)

type RedisClient struct {
	pool   *redis.Pool
	server string
	db     string
}

// The Redis client connection
func NewRedisClient(dsn string) (*RedisClient, error) {
	redisClient := &RedisClient{}
	if err := redisClient.InitDefaultPool(dsn); nil != err {
		return nil, err
	}
	return redisClient, nil
}

func (client *RedisClient) GetServer() string {
	return client.server
}

func (client *RedisClient) GetDb() string {
	return client.db
}

func (client *RedisClient) InitDefaultPool(dsn string) error {
	return client.InitPool(dsn, DefaultMaxIdle, DefaultMaxActive)
}

func (client *RedisClient) InitPool(dsn string, maxIdle, maxActive int) error {
	server, password, db, err := utils.DsnParse(dsn)
	if err != nil {
		return err
	}
	client.server = server
	client.db = db

	client.pool = &redis.Pool{
		MaxIdle:     maxIdle,            // default: 3
		MaxActive:   maxActive,          // default: 1000
		IdleTimeout: DefaultMaxWaitTime, // default 3 * 60 seconds
		Dial: func() (redis.Conn, error) {
			c, err := validateServer(server, password, utils.StringToInt(db))
			if err != nil {
				return nil, err
			}
			log.Println("new redis pool success!")
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return nil
}

func validateServer(server, password string, db int) (redis.Conn, error) {
	c, err := redis.Dial("tcp", server)
	if err != nil {
		log.Println("failed to connect:", err)
		return nil, err
	}
	if password != "" {
		if _, err := c.Do("AUTH", password); err != nil {
			log.Println("password auth fail", err)
			c.Close()
			return nil, err
		}
	}
	if db != DefaultRedisDb {
		if _, err := c.Do("SELECT", db); err != nil {
			c.Close()
			log.Println("select db fail", err)
			return nil, err
		}
	}
	return c, err
}

func (client *RedisClient) GetConnection() (conn redis.Conn) {
	return client.pool.Get()
}

func (client *RedisClient) ReturnConn(conn redis.Conn) {
	if conn != nil {
		err := conn.Close()
		if err != nil {
			log.Fatalln("redis connection close fail", err)
		}
	}
}

func (client *RedisClient) Do(cmd string, args ...interface{}) (reply interface{}, err error) {
	conn := client.GetConnection()
	defer client.ReturnConn(conn)
	return conn.Do(cmd, args...)
}

func (client *RedisClient) Get(key string) (string, error) {
	return redis.String(client.Do("GET", key))
}

func (client *RedisClient) Set(key, value string) (bool, error) {
	res, err := redis.String(client.Do("SET", key, value))
	if err != nil {
		return false, err
	}
	if res == "OK" {
		return true, nil
	} else {
		return false, nil
	}
}

func (client *RedisClient) SetEx(key string, value string, seconds int) (bool, error) {
	resp, err := client.Do("SETEX", key, seconds, value)
	if resp == "OK" {
		return true, nil
	} else {
		return false, err
	}
}

func (client *RedisClient) SetNx(key string, value string, seconds int) (bool, error) {
	res, err := redis.Int(client.Do("SETNX", key, value))
	if err != nil {
		return false, err
	}
	if res == Success {
		return redis.Bool(client.Do("EXPIRE", key, seconds))
	} else {
		return false, nil
	}
}

func (client *RedisClient) Exists(key string) (bool, error) {
	return redis.Bool(client.Do("EXISTS", key))
}

func (client *RedisClient) Del(keys ...string) (bool, error) {
	if len(keys) == 0 {
		return false, errors.New("no keys")
	}
	args := []interface{}{}
	for _, k := range keys {
		args = append(args, k)
	}
	return redis.Bool(client.Do("DEL", args...))
}

func (client *RedisClient) Incr(key string) (bool, error) {
	return redis.Bool(client.Do("INCR", key))
}
