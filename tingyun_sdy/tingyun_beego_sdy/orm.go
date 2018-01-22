package tingyun_beego_sdy

import (
	tingyun "github.com/TingYunAPM/go"
	"github.com/astaxie/beego/orm"
	"github.com/sudiyi/sdy/utils"
)

var dbs = map[string][2]string{}

func RegisterDataBase(aliasName, driverName, dataSource string, params ...int) error {
	host, _, dbName, err := utils.MysqlDsnParse(dataSource)
	if nil != err {
		return err
	}
	if err = orm.RegisterDataBase(aliasName, driverName, dataSource, params...); nil != err {
		return err
	}
	dbs[aliasName] = [2]string{host, dbName}
	return nil
}

type Orm struct {
	orm.Ormer
	host   string
	dbName string
	action *tingyun.Action
	name   string
}

func NewOrm(action *tingyun.Action, name string) *Orm {
	return NewOrmByName(action, name, "default")
}

func NewOrmByName(action *tingyun.Action, name string, ormName string) *Orm {
	o := &Orm{action: action, name: name}
	o.Ormer = orm.NewOrm()
	if "default" != ormName {
		o.Using(ormName)
	}
	if strs, ok := dbs[ormName]; ok {
		o.host, o.dbName = strs[0], strs[1]
	}
	return o
}

func (o *Orm) SetActionAndName(action *tingyun.Action, name string) *Orm {
	o.action = action
	o.name = name
	return o
}

func (o *Orm) QueryM2M(md interface{}, name string) orm.QueryM2Mer {
	return &QueryM2M{o.Ormer.QueryM2M(md, name), o}
}

func (o *Orm) QueryTable(arg interface{}) orm.QuerySeter {
	return &QuerySet{o.Ormer.QueryTable(arg), o}
}

func (o *Orm) createTingyunComponent(op string) *tingyun.Component {
	return o.action.CreateDBComponent(tingyun.ComponentMysql, o.host, o.dbName, "", op, o.name)
}

// CRUDs

func (o *Orm) Read(md interface{}, cols ...string) error {
	component := o.createTingyunComponent("SELECT")
	defer component.Finish()
	return o.Ormer.Read(md, cols...)
}

func (o *Orm) ReadOrCreate(md interface{}, col1 string, cols ...string) (bool, int64, error) {
	component := o.createTingyunComponent("SELECT/INSERT")
	defer component.Finish()
	return o.Ormer.ReadOrCreate(md, col1, cols...)
}

func (o *Orm) Insert(arg interface{}) (int64, error) {
	component := o.createTingyunComponent("INSERT")
	defer component.Finish()
	return o.Ormer.Insert(arg)
}

func (o *Orm) InsertOrUpdate(md interface{}, args ...string) (int64, error) {
	component := o.createTingyunComponent("INSERT/UPDATE")
	defer component.Finish()
	return o.Ormer.InsertOrUpdate(md, args...)
}

func (o *Orm) InsertMulti(bulk int, mds interface{}) (int64, error) {
	component := o.createTingyunComponent("INSERT")
	defer component.Finish()
	return o.Ormer.InsertMulti(bulk, mds)
}

func (o *Orm) Update(md interface{}, cols ...string) (int64, error) {
	component := o.createTingyunComponent("UPDATE")
	defer component.Finish()
	return o.Ormer.Update(md, cols...)
}

func (o *Orm) Delete(md interface{}, cols ...string) (int64, error) {
	component := o.createTingyunComponent("DELETE")
	defer component.Finish()
	return o.Ormer.Delete(md, cols...)
}

type QueryM2M struct {
	orm.QueryM2Mer
	o *Orm
}

func (qm *QueryM2M) Add(args ...interface{}) (int64, error) {
	component := qm.o.createTingyunComponent("INSERT")
	defer component.Finish()
	return qm.QueryM2Mer.Add(args...)
}

func (qm *QueryM2M) Remove(args ...interface{}) (int64, error) {
	component := qm.o.createTingyunComponent("DELETE")
	defer component.Finish()
	return qm.QueryM2Mer.Remove(args...)
}

func (qm *QueryM2M) Exist(arg interface{}) bool {
	component := qm.o.createTingyunComponent("SELECT")
	defer component.Finish()
	return qm.QueryM2Mer.Exist(arg)
}

func (qm *QueryM2M) Clear() (int64, error) {
	component := qm.o.createTingyunComponent("DELETE")
	defer component.Finish()
	return qm.QueryM2Mer.Clear()
}

func (qm *QueryM2M) Count() (int64, error) {
	component := qm.o.createTingyunComponent("SELECT")
	defer component.Finish()
	return qm.QueryM2Mer.Count()
}

type QuerySet struct {
	orm.QuerySeter
	o *Orm
}

func (qs *QuerySet) Count() (int64, error) {
	component := qs.o.createTingyunComponent("SELECT")
	defer component.Finish()
	return qs.QuerySeter.Count()
}

func (qs *QuerySet) Exist() bool {
	component := qs.o.createTingyunComponent("SELECT")
	defer component.Finish()
	return qs.QuerySeter.Exist()
}

func (qs *QuerySet) Update(values orm.Params) (int64, error) {
	component := qs.o.createTingyunComponent("UPDATE")
	defer component.Finish()
	return qs.QuerySeter.Update(values)
}

func (qs *QuerySet) Delete() (int64, error) {
	component := qs.o.createTingyunComponent("DELETE")
	defer component.Finish()
	return qs.QuerySeter.Delete()
}

func (qs *QuerySet) All(container interface{}, cols ...string) (int64, error) {
	component := qs.o.createTingyunComponent("SELECT")
	defer component.Finish()
	return qs.QuerySeter.All(container, cols...)
}

func (qs *QuerySet) One(container interface{}, cols ...string) error {
	component := qs.o.createTingyunComponent("SELECT")
	defer component.Finish()
	return qs.QuerySeter.One(container, cols...)
}

func (qs *QuerySet) Values(results *[]orm.Params, exprs ...string) (int64, error) {
	component := qs.o.createTingyunComponent("SELECT")
	defer component.Finish()
	return qs.QuerySeter.Values(results, exprs...)
}

func (qs *QuerySet) ValuesList(results *[]orm.ParamsList, exprs ...string) (int64, error) {
	component := qs.o.createTingyunComponent("SELECT")
	defer component.Finish()
	return qs.QuerySeter.ValuesList(results, exprs...)
}

func (qs *QuerySet) ValuesFlat(result *orm.ParamsList, expr string) (int64, error) {
	component := qs.o.createTingyunComponent("SELECT")
	defer component.Finish()
	return qs.QuerySeter.ValuesFlat(result, expr)
}

func (qs *QuerySet) RowsToMap(result *orm.Params, keyCol, valueCol string) (int64, error) {
	component := qs.o.createTingyunComponent("SELECT")
	defer component.Finish()
	return qs.QuerySeter.RowsToMap(result, keyCol, valueCol)
}

func (qs *QuerySet) RowsToStruct(ptrStruct interface{}, keyCol, valueCol string) (int64, error) {
	component := qs.o.createTingyunComponent("SELECT")
	defer component.Finish()
	return qs.QuerySeter.RowsToStruct(ptrStruct, keyCol, valueCol)
}
