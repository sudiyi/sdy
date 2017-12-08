package redisclient

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"github.com/sudiyi/sdy/utils"
	"log"
	"sync"
	"time"
)

var redisInstance *RedisClient = nil
var redisOnce sync.Once

type RedisClient struct {
	pool *redis.Pool
}

const (
	Success int = 1 // 成功

	DefaultRedisDb int = 0

	DefaultMaxIdle     int           = 3                 // 空闲连接的最大数目
	DefaultMaxActive   int           = 1000              // 给定时间内最大连接数，为0则连接数没有限制
	DefaultMaxWaitTime time.Duration = 180 * time.Second // Redis最大等待时间
)

// The Redis client connection
func NewRedisClient(dsn string) (*RedisClient, error) {
	pool, err := NewDefaultPool(dsn)
	if err != nil {
		return nil, err
	}
	return &RedisClient{pool: pool}, nil
}

func NewRedisClientOnce(dsn string) (*RedisClient, error) {
	var error error
	redisOnce.Do(func() {
		pool, err := NewDefaultPool(dsn)
		if err != nil {
			error = err
		} else {
			redisInstance = &RedisClient{pool: pool}
		}
	})
	return redisInstance, error
}

func NewDefaultPool(dsn string) (*redis.Pool, error) {
	return NewPool(dsn, DefaultMaxIdle, DefaultMaxActive)
}

func NewPool(dsn string, maxIdle, maxActive int) (*redis.Pool, error) {
	server, password, db, err := utils.DsnParse(dsn)
	database := utils.StringToInt(db)
	if err != nil {
		return nil, err
	}
	return &redis.Pool{
		MaxIdle:     maxIdle,            // default: 3
		MaxActive:   maxActive,          // default: 1000
		IdleTimeout: DefaultMaxWaitTime, // default 3 * 60 seconds
		Dial: func() (redis.Conn, error) {
			c, err := validateServer(server, password, database)
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
	}, nil
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
	return redis.Bool(client.Do("SETEX", key, seconds, value))
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
