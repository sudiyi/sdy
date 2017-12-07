package main

import (
	"fmt"
	"github.com/huhongda/GoToolBox/captcha"
)

func main() {
	cap := captcha.New("redis://localhost:6379", "accessKey", "secretKey", "SMS_***", "signName", true)
	fmt.Println(cap.SmsSend("158***6956"))
}
