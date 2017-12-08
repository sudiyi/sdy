# TingYun Package

## Usage

```
// main.go
func main() {
    tingyun.AppInit("config/tingyun.json")
    defer tingyun.AppStop()
    r := tingyun_gin.Default()
    r.Get('/handle', controller.Handle)    
    r.Run("0.0.0.0:8080")
}

// controller/handle.go

func Handle(c) {
    action := tingyun_gin.FindAction(c)
    component := action.CreateComponent("Handle")
    defer component.Finish()
    services := handle.New(c, action, component).DoSomething()
}

// services/handle/service.go

type HandleService struct {
    ginClient *gin.Context
    action *tingyun.Action
    component *tingyun.Component         
}

func New(c *gin.Context, action *tingyun.Action, component *tingyun.Component) HandleService {
    return &HandleService{ginClient: c, action: action, component: component}
} 

func (h *HandleService) DoSomething() {
    var doSomething = func () {
        h.Do("FirstThing", h.FirstTing())
    }
    h.Do("DoSomething", doSomething)
}

func (h *HandleService) FirstTing() {
    var firstThing = func() {
    
    }
    h.Do("FirstTing", firstThing)
}

func (h *HandleService) Do(methodName string, args...) {
    sub_component := h.component.CreateComponent(methodName)
    defer sub_component.Finish()
    // your logic
}
```