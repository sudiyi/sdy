# utils

## Usage

```
go get github.com/sudiyi/sdy/redis

or 

glide get github.com/sudiyi/sdy/redis
```

```
client := redis.NewRedisClient("redis://localhost:6379/0")
client.Get("redis_key")
client.SetEx("redis_key", "xxxx", 3)
```


## References

[https://github.com/justinyaoqi/redigohelper/blob/master/redigohelper.go](https://github.com/justinyaoqi/redigohelper/blob/master/redigohelper.go)
[https://github.com/ydx00/SenseAdtargeting/blob/9801c938e5215325649553de033654323d981aaf/src/util/RedisClient.go#L113](https://github.com/ydx00/SenseAdtargeting/blob/9801c938e5215325649553de033654323d981aaf/src/util/RedisClient.go#L113)