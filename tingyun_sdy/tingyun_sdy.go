package tingyun_sdy

import (
	"reflect"
	"runtime"

	tingyun "github.com/TingYunAPM/go"
)

func RunWithName(name string, parent *tingyun.Component, action *tingyun.Action, f interface{}, args ...interface{}) {
	vf := reflect.ValueOf(f)
	if reflect.Func != vf.Kind() {
		return
	}

	vfArgs := make([]reflect.Value, 1+len(args))
	for i, arg := range args {
		vfArgs[i+1] = reflect.ValueOf(arg)
	}

	var component *tingyun.Component
	if nil != parent {
		component = parent.CreateComponent(name)
	} else {
		component = action.CreateComponent(name)
	}
	defer component.Finish()

	vfArgs[0] = reflect.ValueOf(component)
	vf.Call(vfArgs)
}

func Run(parent *tingyun.Component, action *tingyun.Action, f interface{}, args ...interface{}) {
	name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	RunWithName(name, parent, action, f, args...)
}
