package gorm_util

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Hitable interface {
	Hit()
}

type GormCtx struct {
	Db      *gorm.DB
	Hitable Hitable
}

func New(dsn string) (*GormCtx, error) {
	g := &GormCtx{}
	var err error
	if g.Db, err = gorm.Open("mysql", dsn); nil != err {
		return nil, err
	}
	g.Db.DB().SetMaxIdleConns(10)
	g.Db.DB().SetMaxOpenConns(100)
	g.Db.LogMode("release" != gin.Mode())
	return g, nil
}

func (g *GormCtx) New() *gorm.DB {
	return g.Db.New()
}

// CRUDs

func (g *GormCtx) outOrPanic(out *gorm.DB) *gorm.DB {
	if nil != out.Error {
		if nil != g.Hitable {
			g.Hitable.Hit()
		}
		panic(out.Error)
	}
	return out
}

func (g *GormCtx) outOrNotfoundOrPanic(out *gorm.DB) *gorm.DB {
	if (nil != out.Error) && !out.RecordNotFound() {
		if nil != g.Hitable {
			g.Hitable.Hit()
		}
		panic(out.Error)
	}
	return out
}

func (g *GormCtx) First(db *gorm.DB, value interface{}, where ...interface{}) *gorm.DB {
	return g.outOrNotfoundOrPanic(db.First(value, where...))
}

func (g *GormCtx) Last(db *gorm.DB, value interface{}, where ...interface{}) *gorm.DB {
	return g.outOrNotfoundOrPanic(db.Last(value, where...))
}

func (g *GormCtx) Find(db *gorm.DB, value interface{}, where ...interface{}) *gorm.DB {
	return g.outOrNotfoundOrPanic(db.Find(value, where...))
}

func (g *GormCtx) Count(db *gorm.DB, value interface{}) *gorm.DB {
	return g.outOrPanic(db.Count(value))
}

func (g *GormCtx) Scan(db *gorm.DB, dest interface{}) *gorm.DB {
	return g.outOrPanic(db.Scan(dest))
}

func (g *GormCtx) Create(db *gorm.DB, value interface{}) *gorm.DB {
	return g.outOrPanic(db.Create(value))
}

func (g *GormCtx) Save(db *gorm.DB, value interface{}) *gorm.DB {
	return g.outOrPanic(db.Save(value))
}

func (g *GormCtx) Update(db *gorm.DB, attrs ...interface{}) *gorm.DB {
	return g.outOrPanic(db.Update(attrs...))
}

func (g *GormCtx) Updates(db *gorm.DB, values interface{}, ignore ...bool) *gorm.DB {
	return g.outOrPanic(db.Updates(values, ignore...))
}

func (g *GormCtx) UpdateColumn(db *gorm.DB, attrs ...interface{}) *gorm.DB {
	return g.outOrPanic(db.UpdateColumn(attrs...))
}

func (g *GormCtx) UpdateColumns(db *gorm.DB, values interface{}) *gorm.DB {
	return g.outOrPanic(db.UpdateColumns(values))
}

func (g *GormCtx) Delete(db *gorm.DB, value interface{}, where ...interface{}) *gorm.DB {
	return g.outOrPanic(db.Delete(value, where...))
}
