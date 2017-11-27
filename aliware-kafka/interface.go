package kafka

// Message MQ message
type Message interface {
	Key() []byte
	Topic() string
	Value() []byte
	Offset() int64
	Partition() int32
}

// Consumer mq consumer client
type Consumer interface {
	Close()
	Messages() <-chan Message
	Errors() <-chan error
	Commit(message Message)
}
