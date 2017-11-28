# aliware kafka

## Usage

```
go get github.com:huhongda/GoToolBox/aliware-kafka

or 

glide get github.com:huhongda/GoToolBox/aliware-kafka
```

* producer


**Sync**

```
import(
    "github.com:huhongda/GoToolBox/aliware-kafka"
)

client := kafka.New(servers, accessKey, password, debug, logger)
producer := client.NewProducer(topic)
producer.produce(key, content)
```


**Async**

```
import(
    "github.com:huhongda/GoToolBox/aliware-kafka"
    "log"
)

var logger = log.New(os.Stderr, "", log.LstdFlags)
client := kafka.New(servers, accessKey, password, debug, logger)
producer := client.NewAsyncProducer(topic)
producer.AsyncProduce(key, content)
```

* consumer
```
import(
    "github.com:huhongda/GoToolBox/aliware-kafka"
    "os"
    "os/signal"
)
var logger = log.New(os.Stderr, "", log.LstdFlags)

signals := make(chan os.Signal, 1)
signal.Notify(signals, os.Interrupt)

client := kafka.New(servers, accessKey, password, debug, logger)
consumer, err := client.NewConsumer(consumerId, topics, offset)
if err != nil {
    panic(err)
}

for {
    select {
    case msg := <-consumer.Messages():
        fmt.Printf(
            "Topic: %s, Key: %s, Partition: %d, Offset: %d, Content: %s \n", msg.Topic(),
            msg.Key(), msg.Partition(), msg.Offset(), string(msg.Value()),
        )
        
        // you can use the msg.Key() for the logic judge and for your business logic
        
        consumer.Commit(msg)
    case <-signals:
        fmt.Println("Stop consumer server...")
        consumer.Close()
        return
    }
}
```

## References

[aliware-kafka-demos](https://github.com/AliwareMQ/aliware-kafka-demos/kafka-go-demo)

[kafka最佳实践](https://help.aliyun.com/document_detail/60691.html?spm=5176.product29530.6.609.FpkKHb)