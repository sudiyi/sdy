## Synopsis

normal component:
```
import (
	"github.com/gin-gonic/gin"
	"github.com/sudiyi/sdy/gin/sdy_tingyun_gin"
	tingyun "github.com/TingYunAPM/go"
)

// Your gin handler in a controller
func Fun(c *gin.Context) {
	// Creates an action wrapping the Tingyun action of the request
	action := sdy_tingyun_gin.NewAction(c)
	out := gin.H{}
	// Runs a Tingyun component in name "root", with function ServeA and its arguments
	action.Run("root", nil, ServeA, 1, out)
	c.JSON(200, out)
}

// Tingyun component function:
// The 1st argument must be an action that you created in handler.
// The 2nd argument must be a component of the caller of this function.
// And the rest are the arguments you passed in action.Run
func ServeA(action *sdy_tingyun_gin.Action, component *tingyun.Component, a int, out gin.H) {
	action.Run("sub", component, ServeB, a, 2, out)
}

func ServeB(action *sdy_tingyun_gin.Action, component *tingyun.Component, a, b int, out gin.H) {
	out["sum"] = a + b
}
```

gorm (mysql) component:
```
import (
	"github.com/sudiyi/sdy/gin/sdy_tingyun_gin"
	tingyun "github.com/TingYunAPM/go"
)

var g *sdy_tingyun_gin.Gorm

func init() {
	// you should create a sdy_tingyun_gin.Gorm instance in a global context
	// NOT in each time requesting
	g, err = sdy_tingyun_gin.NewGorm("user:password@tcp(127.0.0.1:3306)/db")
}

// Your service
func ServeA(action *sdy_tingyun_gin.Action, component *tingyun.Component, a int, out gin.H) {
	db := g.NewDb(action.Action, "your operation name")
	// Now you get a *gorm.DB, use it as usual!
	db.Model(user).UpdateColumn("mobile", "123")
}
```
