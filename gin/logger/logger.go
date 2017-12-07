package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/sudiyi/sdy/utils"
	"os"
	"path"
	"strconv"
	"sync"
)

type fileLogWriter struct {
	sync.RWMutex
	Filename   string
	fileWriter *os.File
	Perm       string
}

func (w *fileLogWriter) StartLogger() *os.File {
	file, _ := w.createLogFile()
	if w.fileWriter != nil {
		w.fileWriter.Close()
	}
	w.fileWriter = file
	return file
}

// Destroy close the file description, close file writer.
func (w *fileLogWriter) Destory() {
	w.fileWriter.Close()
}

// Flush flush file logger.
// there are no buffering messages in file logger in memory.
// flush file means sync file from disk.

func (w *fileLogWriter) Flush() {
	w.fileWriter.Sync()
}

func (w *fileLogWriter) createLogFile() (*os.File, error) {
	// Open the log file
	perm, err := strconv.ParseInt(w.Perm, 8, 64)
	if err != nil {
		return nil, err
	}
	fd, err := os.OpenFile(w.Filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(perm))
	if err == nil {
		// Make sure file perm is user set perm cause of `os.OpenFile` will obey umask
		os.Chmod(w.Filename, os.FileMode(perm))
	}
	return fd, err
}

func Init() *fileLogWriter {
	w := new(fileLogWriter)

	var filename string = "production.log"
	if gin.Mode() == "debug" {
		filename = "development.log"
	}

	if bool, _ := utils.PathExists("logs"); !bool {
		os.MkdirAll("logs", os.ModePerm)
	}
	w.Filename = path.Join("logs", filename)
	w.Perm = "0755"
	return w
}
