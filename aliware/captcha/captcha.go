package captcha

import (
	"encoding/json"
	"errors"
	"github.com/GiterLab/aliyun-sms-go-sdk/sms"
	"github.com/sudiyi/sdy/redisclient"
	"log"
	"strings"
)

const (
	DefaultLen             = 4
	Expiration             = 10 * 60 // captcha expire time
	ExpireLimitEachCaptcha = 60      // twice interval

	DefaultCategory = "sms-code"

	SendSmsOperateTooMuch int = 1
	SendSmsFail           int = 2
	SendSmsSuccess        int = 3
)

type Captcha struct {
	store        *redisclient.RedisClient
	accessKey    string
	secretKey    string
	templateCode string
	signName     string
	debug        bool
	Category     string
}

func New(dsn, accessKey, secretKey, templateCode, signName string, debug bool) *Captcha {
	store, _ := redisclient.NewRedisClient(dsn)
	return &Captcha{
		store:        store,
		accessKey:    accessKey,
		secretKey:    secretKey,
		templateCode: templateCode,
		signName:     signName,
		debug:        debug,
		Category:     DefaultCategory,
	}
}

func (c *Captcha) SmsSend(mobile string) (string, int, error) {
	return c.GenerateAndSend(mobile, Expiration, ExpireLimitEachCaptcha, DefaultLen)
}

func (c *Captcha) SetCategory(category string) *Captcha {
	c.Category = category
	return c
}

func (c *Captcha) SmsVerify(mobile, code string) bool {
	return c.verify(mobile, code)
}

func (c *Captcha) GenerateAndSend(mobile string, ttl, intervalTtl, length int) (string, int, error) {
	intervalKey := c.getIntervalRedisKey(mobile)
	if ok, _ := c.store.Exists(intervalKey); ok {
		return "", SendSmsOperateTooMuch, errors.New("operate too much")
	} else {
		code := string(randStr(length, NUM))
		if ok, err := c.Sending(mobile, map[string]string{"captcha": code}); ok {
			intervalOk, _ := c.store.SetEx(intervalKey, code, intervalTtl)
			ok, _ := c.store.SetNx(c.getRedisKey(mobile), code, ttl)
			if intervalOk && ok {
				return code, SendSmsSuccess, nil
			} else {
				return "", SendSmsFail, err
			}
		} else {
			return "", SendSmsFail, err
		}
	}
}

func (c *Captcha) Sending(mobile string, params map[string]string) (bool, error) {
	sms.HttpDebugEnable = c.debug
	newSms := sms.New(c.accessKey, c.secretKey)
	paramBytes, err := json.Marshal(params)
	if err != nil {
		return false, err
	}
	paramString := string(paramBytes)
	e, err := newSms.SendOne(mobile, c.signName, c.templateCode, paramString)
	if err != nil {
		return false, err
	}
	log.Println("send sms succeed, mobile:", mobile, paramString, e.GetRequestId())
	return true, nil
}

func (c *Captcha) getRedisKey(mobile string) string {
	return strings.Join([]string{mobile, c.Category}, "-")
}

func (c *Captcha) getIntervalRedisKey(mobile string) string {
	return strings.Join([]string{mobile, c.Category, "interval"}, "-")
}

func (c *Captcha) verify(mobile, code string) bool {
	key := c.getRedisKey(mobile)
	realCode, _ := c.store.Get(key)

	if code == realCode {
		if ok, _ := c.store.Del(key); ok {
			return true
		}
	}
	return false
}
