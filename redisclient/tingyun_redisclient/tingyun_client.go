package tingyun_redisclient

import (
	tingyun "github.com/TingYunAPM/go"
	"github.com/sudiyi/sdy/redisclient"
)

type Client struct {
	*redisclient.RedisClient
}

func NewClient(dsn string) (*Client, error) {
	c := &Client{}
	var err error
	if c.RedisClient, err = redisclient.NewRedisClient(dsn); nil != err {
		return nil, err
	}
	return c, nil
}

func (c *Client) Do(action *tingyun.Action, name string, cmd string, args ...interface{}) (reply interface{}, err error) {
	component := action.CreateDBComponent(tingyun.ComponentRedis, c.GetServer(), c.GetDb(), "", cmd, name)
	defer component.Finish()
	return c.RedisClient.Do(cmd, args...)
}
