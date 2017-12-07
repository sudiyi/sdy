## Synopsis
```
import (
	"github.com/sudiyi/sdy/gin/logger_duo"
    "github.com/sudiyi/sdy/gin/middleware/rails_logger"
)

func main() {
    logger := logger_duo.InitLogger("logs/") // 调用者保证目录可用
    logger.LogLevel = logger_v2.NOTICE

    router := gin.New()
    router.Use(rails_logger.RailsLogger())

    logger.Warn("what", 123)
}
