package tingyun_sdy

import (
	"reflect"

	tingyun "github.com/TingYunAPM/go"
)

type Action struct {
	*tingyun.Action
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
