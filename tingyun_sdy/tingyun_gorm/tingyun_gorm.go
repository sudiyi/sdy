package tingyun_gorm

import (
	"runtime"

	tingyun "github.com/TingYunAPM/go"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type Gorm struct {
	DbRoot *gorm.DB
	host   string
	dbName string
}

func NewGorm(dsn string) (*Gorm, error) {
	g := &Gorm{}
	var db *gorm.DB
	my, err := mysql.ParseDSN(dsn)
	if nil != err {
		return nil, err
	}
	g.host, g.dbName = my.Addr, my.DBName

	if db, err = gorm.Open("mysql", dsn); nil != err {
		return nil, err
	}
	g.DbRoot = db
	g.registerTingyun()
	return g, nil
}

func (g *Gorm) registerTingyun() {
	cb := g.DbRoot.Callback()
	cb.Create().Before("gorm:create").Register("ty:before_create", func(scope *gorm.Scope) {
		g.beforeCb(scope, "INSERT")
	})
	cb.Create().After("gorm:create").Register("ty:after_create", func(scope *gorm.Scope) {
		g.afterCb(scope)
	})
	cb.Query().Before("gorm:query").Register("ty:before_query", func(scope *gorm.Scope) {
		g.beforeCb(scope, "SELECT")
	})
	cb.Query().After("gorm:query").Register("ty:after_query", func(scope *gorm.Scope) {
		g.afterCb(scope)
	})
	cb.Update().Before("gorm:update").Register("ty:before_update", func(scope *gorm.Scope) {
		g.beforeCb(scope, "UPDATE")
	})
	cb.Update().After("gorm:update").Register("ty:after_update", func(scope *gorm.Scope) {
		g.afterCb(scope)
	})
	cb.Delete().Before("gorm:delete").Register("ty:before_delete", func(scope *gorm.Scope) {
		g.beforeCb(scope, "DELETE")
	})
	cb.Delete().After("gorm:delete").Register("ty:after_delete", func(scope *gorm.Scope) {
		g.afterCb(scope)
	})
}

func (g *Gorm) NewDbWithName(action *tingyun.Action, name string) *gorm.DB {
	db := g.DbRoot.New()
	db.InstantSet("ty:action", action)
	db.InstantSet("ty:name", name)
	return db
}

func (g *Gorm) NewDb(action *tingyun.Action) *gorm.DB {
	pc, _, _, _ := runtime.Caller(1)
	return g.NewDbWithName(action, runtime.FuncForPC(pc).Name())
}

func (g *Gorm) beforeCb(scope *gorm.Scope, op string) {
	iAction, _ := scope.DB().Get("ty:action")
	if nil == iAction {
		return
	}
	action := iAction.(*tingyun.Action)
	iName, _ := scope.DB().Get("ty:name")
	var name string
	if nil != iName {
		name = iName.(string)
	}

	component := action.CreateDBComponent(tingyun.ComponentMysql, g.host, g.dbName, scope.TableName(), op, name)
	if nil != component {
		scope.InstanceSet("ty:comp", component)
	}
}

func (g *Gorm) afterCb(scope *gorm.Scope) {
	iComponent, _ := scope.InstanceGet("ty:comp")
	if nil != iComponent {
		iComponent.(*tingyun.Component).Finish()
	}
}
