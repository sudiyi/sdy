## Synopsis

main.go
```
package main

import (
	tingyun_beego "github.com/TingYunAPM/go/framework/beego"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sudiyi/sdy/tingyun_sdy/tingyun_beego_sdy"
)

func runBeego() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterModel(&Admins{})
	tingyun_beego_sdy.RegisterDataBase("default", "mysql", "root:password@tcp(127.0.0.1:3306)/test")
	beego.Router("/user", &UserController{})
	tingyun_beego.Run()
}

func main() {
	runBeego()
}
```

controller.go
```
import (
	tingyun_beego "github.com/TingYunAPM/go/framework/beego"
	"github.com/sudiyi/sdy/tingyun_sdy/tingyun_beego_sdy"
)

type UserController struct {
	tingyun_beego.Controller
}

func (c *UserController) Get() {
	action := tingyun_beego_sdy.NewAction(c.Ctx)
	out := map[string]interface{}{}
	action.Run("root", nil, ServeA, 1, out)
	c.Data["json"] = out
	c.ServeJSON()
}
```

service.go
```
import (
	"time"

	tingyun "github.com/TingYunAPM/go"
	"github.com/sudiyi/sdy/tingyun_sdy"
	"github.com/sudiyi/sdy/tingyun_sdy/tingyun_beego_sdy"
)

type Admins struct {
	ID        uint
	Username  string
	Password  string
	Original  bool
	Mobile    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func ServeA(action *tingyun_sdy.Action, component *tingyun.Component, a int, out map[string]interface{}) {
	out["a"] = a + 1
	orm := tingyun_beego_sdy.NewOrm(action.Action, "the name")
	// Now you get a *tingyun_beego_sdy.Orm, use it as beego/orm.Ormer!
	admin := &Admins{Username: "u2", Password: "p2", Mobile: "m2", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	orm.Insert(admin)
}
```
