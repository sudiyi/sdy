package redisclient

import (
	tingyun "github.com/TingYunAPM/go"
)

func (client *RedisClient) TingyunDo(action *tingyun.Action, name string, cmd string, args ...interface{}) (reply interface{}, err error) {
	component := action.CreateDBComponent(tingyun.ComponentRedis, client.server, client.db, "", cmd, name)
	defer component.Finish()
	return client.Do(cmd, args...)
}
