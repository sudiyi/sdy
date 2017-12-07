# Gin Registrar


## Usage

```
import R "github.com/sudiyi/sdy/gin/registrar"

router := gin.Default()
R.Register(router, "debug", "noroute", "monitors")
```

### Debug
    provide the debug interface, include goroutines count & debug profile
    
    `github.com/gin-contrib/pprof`
    
### Noroute
    provide the no route, return the json data for 404 

### Monitors
    provide the monitors interface, head request `/monitors/status`

