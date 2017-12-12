package sdy_tingyun_gin

import (
	tingyun "github.com/TingYunAPM/go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sudiyi/sdy/utils"
)

type Gorm struct {
	Db     *gorm.DB
	host   string
	dbName string
}

func New(dsn string) (*Gorm, error) {
	g := &Gorm{}
	var err error
	if g.host, _, g.dbName, err = utils.DsnParse(dsn); nil != err {
		return nil, err
	}

	if g.Db, err = gorm.Open("mysql", dsn); nil != err {
		return nil, err
	}
	g.Db.DB().SetMaxIdleConns(10)
	g.Db.DB().SetMaxOpenConns(100)
	g.Db.LogMode("release" != gin.Mode())
	return g, nil
}

func (g *Gorm) New() *gorm.DB {
	return g.Db.New()
}

func (g *Gorm) createTingyunComponent(action *tingyun.Action, table string, op string, name string) *tingyun.Component {
	return action.CreateDBComponent(tingyun.ComponentMysql, g.host, g.dbName, table, op, name)
}

func (g *Gorm) First(action *tingyun.Action, table string, name string, db *gorm.DB, value interface{}, where ...interface{}) *gorm.DB {
	component := g.createTingyunComponent(action, table, "SELECT", name)
	defer component.Finish()
	return db.First(value, where...)
}

func (g *Gorm) Last(action *tingyun.Action, table string, name string, db *gorm.DB, value interface{}, where ...interface{}) *gorm.DB {
	component := g.createTingyunComponent(action, table, "SELECT", name)
	defer component.Finish()
	return db.Last(value, where...)
}

func (g *Gorm) Find(action *tingyun.Action, table string, name string, db *gorm.DB, value interface{}, where ...interface{}) *gorm.DB {
	component := g.createTingyunComponent(action, table, "SELECT", name)
	defer component.Finish()
	return db.Find(value, where...)
}

func (g *Gorm) Count(action *tingyun.Action, table string, name string, db *gorm.DB, value interface{}) *gorm.DB {
	component := g.createTingyunComponent(action, table, "SELECT", name)
	defer component.Finish()
	return db.Count(value)
}

func (g *Gorm) Scan(action *tingyun.Action, table string, name string, db *gorm.DB, dest interface{}) *gorm.DB {
	component := g.createTingyunComponent(action, table, "SELECT", name)
	defer component.Finish()
	return db.Scan(dest)
}

func (g *Gorm) Create(action *tingyun.Action, table string, name string, db *gorm.DB, value interface{}) *gorm.DB {
	component := g.createTingyunComponent(action, table, "INSERT", name)
	defer component.Finish()
	return db.Create(value)
}

func (g *Gorm) Save(action *tingyun.Action, table string, name string, db *gorm.DB, value interface{}) *gorm.DB {
	component := g.createTingyunComponent(action, table, "INSERT/UPDATE", name)
	defer component.Finish()
	return db.Save(value)
}

func (g *Gorm) Update(action *tingyun.Action, table string, name string, db *gorm.DB, attrs ...interface{}) *gorm.DB {
	component := g.createTingyunComponent(action, table, "UPDATE", name)
	defer component.Finish()
	return db.Update(attrs...)
}

func (g *Gorm) Updates(action *tingyun.Action, table string, name string, db *gorm.DB, values interface{}, ignore ...bool) *gorm.DB {
	component := g.createTingyunComponent(action, table, "UPDATE", name)
	defer component.Finish()
	return db.Updates(values, ignore...)
}

func (g *Gorm) UpdateColumn(action *tingyun.Action, table string, name string, db *gorm.DB, attrs ...interface{}) *gorm.DB {
	component := g.createTingyunComponent(action, table, "UPDATE", name)
	defer component.Finish()
	return db.UpdateColumn(attrs...)
}

func (g *Gorm) UpdateColumns(action *tingyun.Action, table string, name string, db *gorm.DB, values interface{}) *gorm.DB {
	component := g.createTingyunComponent(action, table, "UPDATE", name)
	defer component.Finish()
	return db.UpdateColumns(values)
}

func (g *Gorm) Delete(action *tingyun.Action, table string, name string, db *gorm.DB, value interface{}, where ...interface{}) *gorm.DB {
	component := g.createTingyunComponent(action, table, "DELETE", name)
	defer component.Finish()
	return db.Delete(value, where...)
}
