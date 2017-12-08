## Synopsis
```
import (
    "github.com/sudiyi/sdy/gin/logger_duo"
    "github.com/sudiyi/sdy/gin/middleware/rails_logger"
)

func main() {
    logger := logger_duo.InitLogger("logs/")
    logger.LogLevel = logger_duo.NOTICE

    router := gin.New()
    router.Use(rails_logger.RailsLogger())

    logger.Warn("what", 123)
}
