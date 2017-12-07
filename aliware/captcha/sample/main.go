package main

import (
	"fmt"
	"github.com/sudiyi/sdy/captcha"
)

func main() {
	cap := captcha.New("redis://localhost:6379", "LTAIKgSfnWzmXz1c", "DcCbEjzNjv9GYABtPbhs31CictXZXD", "SMS_44290055", "58分享", true)
	fmt.Println(cap.SmsSend("15828566956"))
}
