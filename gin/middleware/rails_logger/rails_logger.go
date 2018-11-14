package rails_logger

import (
	"bytes"
	"fmt"
	"github.com/satori/go.uuid"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

type writer struct {
	gin.ResponseWriter
	dup *bytes.Buffer
}

func (w *writer) Write(data []byte) (int, error) {
	if w.Status() >= http.StatusBadRequest {
		w.dup.Write(data)
	}
	return w.ResponseWriter.Write(data)
}

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return green
	case code >= 300 && code < 400:
		return white
	case code >= 400 && code < 500:
		return yellow
	default:
		return red
	}
}

func colorForMethod(method string) string {
	switch method {
	case "GET":
		return blue
	case "POST":
		return cyan
	case "PUT":
		return yellow
	case "DELETE":
		return red
	case "PATCH":
		return green
	case "HEAD":
		return magenta
	case "OPTIONS":
		return white
	default:
		return reset
	}
}

func RailsLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		log(c, gin.DefaultWriter)
	}
}

func log(c *gin.Context, out io.Writer) {
	start := time.Now()
	dup := &bytes.Buffer{}
	c.Writer = &writer{c.Writer, dup}
	reqBody, _ := ioutil.ReadAll(c.Request.Body)
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
	c.Set("request_id", getRequestID())

	c.Next()

	end := time.Now()
	statusCode := c.Writer.Status()
	method := c.Request.Method
	url := c.Request.URL

	respBody := make([]byte, 128)
	respBodyLen, _ := dup.Read(respBody)
	fmt.Fprintf(out, "[G] %11v [%v] %15s %s %s%3d%s %s%s%s %s, Query: %s Parameters: %s | %s\n",
		end.Sub(start),
		end.Format("2006/01/02 15:04:05"),
		c.ClientIP(),
		c.GetString("request_id"),
		colorForStatus(statusCode), statusCode, reset,
		colorForMethod(method), method, reset,
		url.Path,
		url.RawQuery,
		reqBody,
		string(respBody[:respBodyLen]),
	)
}

func getRequestID() string {
	id,_:=uuid.NewV4()
	return id.String()
}