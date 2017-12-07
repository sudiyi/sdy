package logger_duo

import (
	"log"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/sudiyi/sdy/utils"
)

func InitLogger(logPath string) (*Logger, error) {
	if ok, _ := utils.PathExists(logPath); !ok {
		if err := os.MkdirAll(logPath, 0777); nil != err {
			return nil, err
		}
	}

	l, err := Init(&InitArgs{
		Filename: path.Join(logPath, gin.Mode()+".error.log"),
		Flags:    log.LstdFlags | log.Lshortfile,
	})
	if nil != err {
		return nil, err
	}

	la, err := Init(&InitArgs{
		Filename: path.Join(logPath, gin.Mode()+".access.log"),
	})
	if nil != err {
		return nil, err
	}
	gin.DefaultWriter = la.GetWriter()
	return l, nil
}
