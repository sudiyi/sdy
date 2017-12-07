# aliware kafka

## Usage

```
go get github.com/sudiyi/sdy/aliware/kafka

or 

glide get github.com/sudiyi/sdy/aliware/kafka
```

* producer


**Sync**

```
import(
    "github.com/sudiyi/sdy/aliware/kafka"
)

client := kafka.New(servers, accessKey, password, debug)
producer, err := client.NewProducer(topic)
if err != nil {
    panic(err)
}
producer.produce(key, content)
```


**Async**

```
import(
    "github.com/sudiyi/sdy/aliware/kafka"
    "log"
)

client := kafka.New(servers, accessKey, password, debug)
producer, err := client.NewAsyncProducer(topic)
if err != nil {
    panic(err)
}
producer.AsyncProduce(key, content)
```

* consumer
```
import(
    "github.com/sudiyi/sdy/aliware/kafka"
    "os"
    "os/signal"
)

signals := make(chan os.Signal, 1)
signal.Notify(signals, os.Interrupt)

client := kafka.New(servers, accessKey, password, debug)
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