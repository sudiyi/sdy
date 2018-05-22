package kafka

import (
	"fmt"
	"github.com/satori/go.uuid"
	"time"
)

type Wrapper struct {
	Header HeaderWrapper          `json:"header"`
	Body   map[string]interface{} `json:"body"`
}

type HeaderWrapper struct {
	Timestamp int64  `json:"timestamp"`
	Type      int    `json:"type"`
	Uid       string `json:"uid"`
	Timeout   int32  `json:"timeout"`
	UserAgent string `json:"user_agent"`
}

func ContentWrapper(kafkaType int, userAgent string, content map[string]interface{}) *Wrapper {
	return &Wrapper{
		Header: *headerWrapper(kafkaType, userAgent),
		Body:   content,
	}
}

func headerWrapper(kafkaType int, userAgent string) *HeaderWrapper {
	result := uuid.NewV4()
	uid := fmt.Sprintf("%s", result)
	return &HeaderWrapper{
		Timestamp: time.Now().UnixNano() / 1000 / 1000,
		Type:      kafkaType,
		Uid:       uid,
		UserAgent: userAgent,
		Timeout:   3000,
	}
}
