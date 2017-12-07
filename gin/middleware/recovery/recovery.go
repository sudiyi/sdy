package recovery

import (
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/sudiyi/sdy/gin/logger_duo"
)

func recovery(c *gin.Context, logger *logger_duo.Logger, callback func(p interface{})) {
	p := recover()
	if nil == p {
		return
	}

	if nil != callback {
		callback(p)
	}
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file, line = "?", 0
	}
	logger.Errf("panic on %s:%d:", file, line)
	logger.Err(p)
	c.AbortWithStatus(500)
}

func Recovery(logger *logger_duo.Logger, callback func(p interface{})) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer recovery(c, logger, callback)
		c.Next()
	}
}
