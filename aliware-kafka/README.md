# 阿里云kafka

## Usage

```
go get github.com:huhongda/GoToolBox/aliware-kafka

or 

glide get github.com:huhongda/GoToolBox/aliware-kafka
```

* producer
```
import(
    "github.com:huhongda/GoToolBox/aliware-kafka"
)

client := kafka.New(servers, accessKey, password, debug, logger)
producer := client.NewProducer(topic)
producer.produce(key, content)
```

* consumer
```
import(
    "github.com:huhongda/GoToolBox/aliware-kafka"
)

client := kafka.New(servers, accessKey, password, debug, logger)
consumer := client.NewConsumer(consumerId, topics)

channel := consumer.consume()
for val := range channel {
    fmt.Println("out: %s", val)
    //c <- ""
} 
```

## References

[aliware-kafka-demos](https://github.com/AliwareMQ/aliware-kafka-demos/kafka-go-demo)

[kafka最佳实践](https://help.aliyun.com/document_detail/60691.html?spm=5176.product29530.6.609.FpkKHb)