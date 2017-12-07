# Gin Logger to file

## Usage

```
import "github.com/sudiyi/sdy/gin/logger"


logger := Logger.Init()
gin.DefaultWriter = logger.StartLogger()

```

NOTES: the code must before `router := gin.New()`

## NOTES

this package will create a `logs` directory, and the logger will input to `development.log` or `production.log`

## Demo

```
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sudiyi/sdy/gin/logger"
)

func main() {
	logger := logger.Init()
	gin.DefaultWriter = logger.StartLogger()
    
    router := gin.New()
    router.Use(gin.Logger())
    router.Use(gin.Recovery())

	router.GET("/helloworld", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"hello": "world"})
	})
	router.Run(":8080")
}
```