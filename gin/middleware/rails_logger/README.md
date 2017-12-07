# RailsLogger Logger Format

## Usage

```

package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/sudiyi/sdy/gin/logger"
	. "github.com/sudiyi/sdy/gin/middleware"
	R "github.com/sudiyi/sdy/gin/registrar"
)

func InitRouters() *gin.Engine {
	logger := logger.Init()
	gin.DefaultWriter = logger.StartLogger()

	router := gin.New()         // NOTES: please use gin.New() replace the gin.Default()
	router.Use(RailsLogger())   // add RailsLogger() middleware to gin
	router.Use(gin.Recovery())

	router = SetBmsRouter(router)
	R.Register(router, "debug", "noroute", "monitors")
	return router
}
```

## Logger Format

```
// 200 Response
[G] 34.996862ms [2017/08/17 16:50:31]       127.0.0.1 200 POST /bms/v1/devices/configurations, Parameters: data=%7B%22device_id%22%3A%2229366246%22%7D&sign=5eb215cae6346bdeeb4fadccde148a1c&business_type=1&device_id=29366246 |

// > 400 Response will record the response body
[G] 27.639877ms [2017/08/17 16:53:13]       127.0.0.1 422 POST /bms/v1/devices/configurations, Parameters: data=%7B%22device_id%22%3A%2229366246%22%7D&sign=5eb215cae6346bdeeb4fadccde148a1c&business_type=1&device_id=29366246 | {"code":1,"data":{},"message":"device not found"}
```