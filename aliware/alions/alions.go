package alions

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/go-resty/resty"
	"github.com/sudiyi/sdy/utils"
	"net/url"
	"strconv"
	"strings"
	"time"
	"bytes"
)

type AlionsConfigs struct {
	// 消息主体
	Topic string `ini:"topic"`
	// 公测URL
	Url string `ini:"url"`
	// 访问码
	UserAccessKey string `ini:"user_access_key"`
	// 密钥
	UserSecretKey string `ini:"user_secret_key"`
	// 生产者ID
	ProducerGroup  string `ini:"producer_group"`
	Tag            string `ini:"tag"`
	WriteTimeout   int    `ini:"write_timeout"`
	ConnectTimeout int    `ini:"connect_timeout"`
	ReadTimeout    int    `ini:"read_timeout"`
	ConsumerGroup  string `ini:"consumer_group"`
	HeaderTimeout  int    `ini:"alions_header_timeout"`
}

var AlionsConfig *AlionsConfigs

func init() {
	if AlionsConfig == nil {
		AlionsConfig = &AlionsConfigs{
			Topic:          "TEST",
			Url:            "http://test_uri",
			UserAccessKey:  "asdfas",
			UserSecretKey:  "asdfasdfs",
			ProducerGroup:  "PID-TERMINAL-MQ",
			ConsumerGroup:  "xxx",
			Tag:            "DOWN-MESSAGE-V2",
			WriteTimeout:   2,
			ConnectTimeout: 15,
			ReadTimeout:    10,
			HeaderTimeout:  5000,
		}
	}
}

func CurrentTimeForMillisSecond() int64 {
	return time.Now().UnixNano() / 1000000
}

/*
	go 获取openssl hmac sha1 的base64
*/
func sign(data string) string {
	secretKey := AlionsConfig.UserSecretKey
	keyForSign := []byte(secretKey)
	h := hmac.New(sha1.New, keyForSign)
	h.Write([]byte(data))
	encodedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return strings.TrimSpace(encodedStr)
}

func postMd5Str(body string) string {
	data := []byte(body)
	return fmt.Sprintf("%x", md5.Sum(data))
}

func getSignStr(httpMethod string, time int64, str string) string {
	resultStr := bytes.NewBufferString(AlionsConfig.Topic + "\n")
	if httpMethod == "POST" {
		resultStr.WriteString(AlionsConfig.ProducerGroup + "\n")
	 	resultStr.WriteString(postMd5Str(str) + "\n")
	} else {
		resultStr.WriteString(AlionsConfig.ConsumerGroup + "\n")
		if str != "" {
			resultStr.WriteString(str + "\n")
		}
	}
	resultStr.WriteString(strconv.FormatInt(time, 10))
	return sign(resultStr.String())
}

func getHeader(httpMethod string, time int64, str string) map[string]string {
	header := map[string]string{
		"AccessKey": AlionsConfig.UserAccessKey,
		"Signature": getSignStr(httpMethod, time, str),
	}
	if httpMethod == "POST" {
		header["ProducerId"] = AlionsConfig.ProducerGroup
	} else {
		header["ConsumerId"] = AlionsConfig.ConsumerGroup
	}
	return header
}

func request(method, data string, time int64, queryParams map[string]string) (responseBody map[string]interface{}) {
	var resp *resty.Response
	var err error
	switch method {
	case "GET":
		resp, err = resty.R().
			SetHeaders(getHeader("GET", time, data)).
			SetQueryParams(queryParams).
			Get(AlionsConfig.Url)
	case "POST":
		resp, err = resty.R().
			SetHeaders(getHeader("POST", time, data)).
			SetBody(data).
			SetQueryParams(queryParams).
			Post(AlionsConfig.Url)
	case "DELETE":
		resp, err = resty.R().SetHeaders(getHeader("DELETE", time, data)).Delete(AlionsConfig.Url)
	}
	if err != nil {
		panic(err)
	}
	body := resp.Body()
	responseBody, _ = utils.JsonToMap(string(body))
	responseBody["code"] = resp.StatusCode()
	return responseBody
}

func Post(body string, tag string, key string) (responseBody map[string]interface{}) {
	time := CurrentTimeForMillisSecond()
	_body := url.QueryEscape(body)
	if tag == "" {
		tag = "http"
	}

	if key == "" {
		key = "http"
	}

	args := map[string]string{
		"topic": AlionsConfig.Topic,
		"time":  strconv.FormatInt(time, 10),
		"tag":   tag,
		"key":   key,
	}
	return request("POST", _body, time, args)
}

func Get() (responseBody map[string]interface{}) {
	time := CurrentTimeForMillisSecond()
	args := map[string]string{
		"topic": AlionsConfig.Topic,
		"time":  strconv.FormatInt(time, 10),
		"num":   "32",
	}

	return request("GET", "", time, args)
}

func Delete(msg_handle string) (responseBody map[string]interface{}) {
	time := CurrentTimeForMillisSecond()
	return request("DELETE", msg_handle, time, nil)
}
