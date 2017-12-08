package monitor

import (
	"errors"
	"github.com/TingYunAPM/go"
	"github.com/TingYunAPM/go/framework/gin"
	"github.com/gin-gonic/gin"
	"reflect"
	"sdy/utils"
)

type TingYunService struct {
	action    *tingyun.Action
	component *tingyun.Component
}

// use for main init method: monitor.AppInit()
func AppInit(file ...string) {
	var configPath string
	if len(file) == 0 {
		configPath = "config/tingyun.json"
	} else {
		configPath = file[len(file)-1]
	}
	tingyun.AppInit(configPath)
}

// use for main method: defer monitor.AppStop()
func AppStop() {
	tingyun.AppStop()
}

// router := monitor.GinDefault()
func GinDefault() *tingyun_gin.WrapEngine {
	return tingyun_gin.Default()
}

// monitor := monitor.New("handle", c)
func New(handleName string, c *gin.Context) *TingYunService {
	action := tingyun_gin.FindAction(c)
	component := action.CreateComponent(handleName)
	return &TingYunService{action: action, component: component}
}

func (t *TingYunService) EnableRedis(dsn string) *tingyun.Component {
	host, _, db, _ := utils.DsnParse(dsn)
	return t.action.CreateDBComponent(tingyun.ComponentRedis, host, db, "", "get/set/de", "redis.Do")
}

func (t *TingYunService) WrapperRun(m map[string]interface{}, name string, params ...interface{}) (result []reflect.Value, err error) {
	subComponent := t.component.CreateComponent(name)
	defer subComponent.Finish()

	f := reflect.ValueOf(m[name])
	if len(params) != f.Type().NumIn() {
		err = errors.New("The number of params is not adapted.")
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	result = f.Call(in)
	return
}
