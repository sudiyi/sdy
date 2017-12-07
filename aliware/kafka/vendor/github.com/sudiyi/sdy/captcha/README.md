# sms-captcha

Golang sms captcha, special for aliYun sms

## dependences

* [redigo](https://github.com/garyburd/redigo)
* [aliyun-sms-go-sdk](https://github.com/GiterLab/aliyun-sms-go-sdk)

## Usage

```
go get github.com/huhongda/sms-captcha

or 

glide get github.com/huhongda/sms-captcha
```

```
redis+sentinel://[:password@]host:port[,host2:port2,...][/service_name[/db]][?param1=value1[&param2=value=2&...]]
redis_dsn = "redis://:password@localost:6379/10"
```

```
import (
    captcha "github.com/huhongda/sms-captcha"
) 

cap = captcha.New("redis://:password@localost:6379/10", "accessKey", "secretKey", "SMS_******", "签名", false)
cap.SmsSend("15****6956")
```


### Future

```
cap = captcha.New("redis://:password@localost:6379/10", "accessKey", "secretKey", "SMS_******", "签名", false)
cap.SetCategroy("sms-code") // your can give the default value
cap.SmsSend("15****6956")
```