package tingyun_redisclient

import (
	tingyun "github.com/TingYunAPM/go"
	"github.com/sudiyi/sdy/redisclient"
)

func TingAndDo(action *tingyun.Action, name string, red *redisclient.RedisClient) func(string, ...interface{}) (interface{}, error) {
	return func(cmd string, args ...interface{}) (interface{}, error) {
		component := action.CreateDBComponent(tingyun.ComponentRedis, red.GetServer(), red.GetDb(), "", cmd, name)
		defer component.Finish()
		return red.Do(cmd, args...)
	}
}
