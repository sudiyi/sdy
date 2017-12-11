package sdy_tingyun_gin

import (
	"reflect"

	tingyun "github.com/TingYunAPM/go"
	tingyun_gin "github.com/TingYunAPM/go/framework/gin"
	"github.com/gin-gonic/gin"
)

type Action struct {
	*tingyun.Action
}

func NewAction(c *gin.Context) *Action {
	return &Action{tingyun_gin.FindAction(c)}
}

const argsBeginIndex = 2

func (a *Action) Run(name string, parent *tingyun.Component, f interface{}, args ...interface{}) {
	vf := reflect.ValueOf(f)
	if reflect.Func != vf.Kind() {
		return
	}

	vfArgs := make([]reflect.Value, argsBeginIndex+len(args))
	vfArgs[0] = reflect.ValueOf(a)
	for i, arg := range args {
		vfArgs[argsBeginIndex+i] = reflect.ValueOf(arg)
	}

	var component *tingyun.Component
	if nil != parent {
		component = parent.CreateComponent(name)
	} else {
		component = a.CreateComponent(name)
	}
	defer component.Finish()

	vfArgs[1] = reflect.ValueOf(component)
	vf.Call(vfArgs)
}
