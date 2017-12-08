package main

import (
	"errors"
	"fmt"
	"reflect"
)

func main() {
	a := map[string]interface{}{"a": func() (int, error) {
		fmt.Println("1111")
		return 1111, errors.New("test")
	}}
	value, err := Call(a, "a")
	fmt.Println(value)
	fmt.Println(err)
}

func Call(m map[string]interface{}, name string, params ...interface{}) (result []reflect.Value, err error) {
	f := reflect.ValueOf(m[name])
	if len(params) != f.Type().NumIn() {
		err = errors.New("The number of params is not adapted.")
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	result = f.Call(in)
	return
}
