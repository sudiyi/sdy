## Synopsis

```
import (
	tingyun "github.com/TingYunAPM/go"
	tingyun_gin "github.com/TingYunAPM/go/framework/gin"
	"github.com/gin-gonic/gin"
	"github.com/sudiyi/sdy/tingyun_sdy"
	"github.com/sudiyi/sdy/tingyun_sdy/tingyun_gin_sdy"
)

func main() {
	r := tingyun_gin.New()
	r.GET("/f", Fun)
}

// Your gin handler in a controller
func Fun(c *gin.Context) {
	// Creates an action wrapping the Tingyun action of the request
	action := tingyun_gin_sdy.NewAction(c)
	out := gin.H{}
	// Runs a Tingyun component in name "root", with function ServeA and its arguments
	action.Run("root", nil, ServeA, 1, out)
	c.JSON(200, out)
}

// Tingyun component function:
// The 1st argument must be an action that you created in handler.
// The 2nd argument must be a component of the caller of this function.
// And the rest are the arguments you passed in action.Run
func ServeA(action *tingyun_sdy.Action, component *tingyun.Component, a int, out gin.H) {
	action.Run("sub", component, ServeB, a, 2, out)
}

func ServeB(action *tingyun_sdy.Action, component *tingyun.Component, a, b int, out gin.H) {
	out["sum"] = a + b
}
```
