package tingyun_beego_sdy

import (
	tingyun_beego "github.com/TingYunAPM/go/framework/beego"
	"github.com/astaxie/beego/context"
	"github.com/sudiyi/sdy/tingyun_sdy"
)

func NewAction(ctx *context.Context) *tingyun_sdy.Action {
	return &tingyun_sdy.Action{tingyun_beego.FindAction(ctx)}
}
