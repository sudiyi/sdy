package basic_redis_auth

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

const RED_KEY_PREF = "bauth:"

type BasicAuth struct {
	red    *redis.Client
	prefix string
	realm  string
}

func New(red *redis.Client, namespace, realm string) *BasicAuth {
	prefix := RED_KEY_PREF
	if len(namespace) > 0 {
		prefix = namespace + ":" + prefix
	}
	return &BasicAuth{
		red:    red,
		prefix: prefix,
		realm:  realm,
	}
}

func decodeAuth(s string) (string, string, bool) {
	if !strings.HasPrefix(s, "Basic ") {
		return "", "", false
	}
	s = s[6:]
	if 0 == len(s) {
		return "", "", false
	}

	bs, err := base64.StdEncoding.DecodeString(s)
	if nil != err {
		return "", "", false
	}
	s = string(bs)

	es := strings.SplitN(s, ":", 2)
	if 2 != len(es) {
		return "", "", false
	}
	return es[0], es[1], true
}

func (b *BasicAuth) auth(c *gin.Context) int {
	u, p, ok := decodeAuth(c.GetHeader("Authorization"))
	if !ok {
		return http.StatusUnauthorized
	}

	rp, err := b.red.Get(b.prefix + u).Result()
	if redis.Nil == err {
		return http.StatusUnauthorized
	} else if nil != err {
		return http.StatusInternalServerError
	} else if rp != p {
		return http.StatusUnauthorized
	}
	return http.StatusOK
}

func (b *BasicAuth) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		status := b.auth(c)
		if http.StatusOK != status {
			if http.StatusUnauthorized == status {
				c.Header("WWW-Authenticate", fmt.Sprintf(`Basic Realm="%s"`, b.realm))
			}
			c.AbortWithStatus(status)
			return
		}
		c.Next()
	}
}
