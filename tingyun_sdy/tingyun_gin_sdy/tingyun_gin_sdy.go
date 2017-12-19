package tingyun_gin_sdy

import (
	tingyun_gin "github.com/TingYunAPM/go/framework/gin"
	"github.com/gin-gonic/gin"
	"github.com/sudiyi/sdy/tingyun_sdy"
)

func NewAction(c *gin.Context) *tingyun_sdy.Action {
	return &tingyun_sdy.Action{tingyun_gin.FindAction(c)}
}
