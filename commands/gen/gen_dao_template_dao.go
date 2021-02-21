package gen

const templateDaoDaoIndexContent = `
// ============================================================================
// This is auto-generated by gf cli tool only once. Fill this file as you wish.
// ============================================================================

package dao

import (
	"{TplImportPrefix}/dao/internal"
)

// {TplTableNameCamelLowerCase}Dao is the manager for logic model data accessing
// and custom defined data operations functions management. You can define
// methods on it to extend its functionality as you wish.
type {TplTableNameCamelLowerCase}Dao struct {
	*internal.{TplTableNameCamelCase}Dao
}

var (
	// {TplTableNameCamelCase} is globally public accessible object for table {TplTableName} operations.
	{TplTableNameCamelCase} = &{TplTableNameCamelLowerCase}Dao{
		internal.{TplTableNameCamelCase},
	}
)

// Fill with you ideas below.

`

const templateDaoDaoInternalContent = `
// ==========================================================================
// 这是由gf cli工具自动生成的。不要手动编辑此文件。
// ==========================================================================

package internal

import (
	"context"
	"database/sql"
	"{TplImportPrefix}/model"
	"github.com/kotlin2018/orm"
	"github.com/kotlin2018/pkg/g"
	"github.com/kotlin2018/pkg/gmvc"
	"time"
)

// {TplTableNameCamelCase}Dao 是用于逻辑模型数据访问和自定义数据操作功能管理的管理器
type {TplTableNameCamelCase}Dao struct {
	gmvc.M
	DB      orm.DB
	Table   string
	Columns {TplTableNameCamelLowerCase}Columns
}

// {TplTableNameCamelCase}列定义并存储表{TplTableName}的列名
type {TplTableNameCamelLowerCase}Columns struct {
	{TplColumnDefine}
}

var (
	// {TplTableNameCamelCase} 是表{TplTableName}操作的全局公共可访问对象
	{TplTableNameCamelCase} = &{TplTableNameCamelCase}Dao{
		M:     g.DB("{TplGroupName}").Model("{TplTableName}").Safe(),
		DB:    g.DB("{TplGroupName}"),
		Table: "{TplTableName}",
		Columns: {TplTableNameCamelLowerCase}Columns{
			{TplColumnNames}
		},
	}
)

// Ctx 是一个链接函数，它创建并返回一个新的DB，该DB是当前DB对象的浅层副本，其中包含给定的上下文。
func (d *{TplTableNameCamelCase}Dao) Ctx(ctx context.Context) *{TplTableNameCamelCase}Dao {
	return &{TplTableNameCamelCase}Dao{M: d.M.Ctx(ctx)}
}

// As 设置当前表的别名。
func (d *{TplTableNameCamelCase}Dao) As(as string) *{TplTableNameCamelCase}Dao {
	return &{TplTableNameCamelCase}Dao{M: d.M.As(as)}
}

// TX 设置当前操作的事务。
func (d *{TplTableNameCamelCase}Dao) TX(tx *orm.TX) *{TplTableNameCamelCase}Dao {
	return &{TplTableNameCamelCase}Dao{M: d.M.TX(tx)}
}

// Master 指定操作是在主节点上进行。
func (d *{TplTableNameCamelCase}Dao) Master() *{TplTableNameCamelCase}Dao {
	return &{TplTableNameCamelCase}Dao{M: d.M.Master()}
}

// Slave 指定操作是在从节点上执行。(请注意，只有在配置了任何从属节点时才有意义)
func (d *{TplTableNameCamelCase}Dao) Slave() *{TplTableNameCamelCase}Dao {
	return &{TplTableNameCamelCase}Dao{M: d.M.Slave()}
}

// Args 为模型操作设置自定义参数。
func (d *{TplTableNameCamelCase}Dao) Args(args ...interface{}) *{TplTableNameCamelCase}Dao {
	return &{TplTableNameCamelCase}Dao{M: d.M.Args(args ...)}
}

// LeftJoin 对模型执行“left join ... on ...”语句。
// 参数<table>可以是联接表及其联接条件，也可以是其别名，例如:
// Table("user").LeftJoin("user_detail", "user_detail.uid=user.uid")
// Table("user", "u").LeftJoin("user_detail", "ud", "ud.uid=u.uid")
func (d *{TplTableNameCamelCase}Dao) LeftJoin(table ...string) *{TplTableNameCamelCase}Dao {
	return &{TplTableNameCamelCase}Dao{M: d.M.LeftJoin(table...)}
}

// RightJoin 对模型执行“right join ... on ...”语句。
// 参数<table>可以是联接表及其联接条件，也可以是其别名，例如:
// Table("user").RightJoin("user_detail", "user_detail.uid=user.uid")
// Table("user", "u").RightJoin("user_detail", "ud", "ud.uid=u.uid")
func (d *{TplTableNameCamelCase}Dao) RightJoin(table ...string) *{TplTableNameCamelCase}Dao {
	return &{TplTableNameCamelCase}Dao{M: d.M.RightJoin(table...)}
}

// InnerJoin 对模型执行“inner join ... on ...”语句。
// 参数<table>可以是联接表及其联接条件，也可以是其别名，例如:
// Table("user").InnerJoin("user_detail", "user_detail.uid=user.uid")
// Table("user", "u").InnerJoin("user_detail", "ud", "ud.uid=u.uid")
func (d *{TplTableNameCamelCase}Dao) InnerJoin(table ...string) *{TplTableNameCamelCase}Dao {
	return &{TplTableNameCamelCase}Dao{M: d.M.InnerJoin(table...)}
}

// Fields 指定需要操作的表字段，多个字段使用字符','连接。
// 参数<fieldNamesOrMapStruct>的类型可以是string/map/*map/struct/*struct
func (d *{TplTableNameCamelCase}Dao) Fields(fieldNamesOrMapStruct ...interface{}) *{TplTableNameCamelCase}Dao {
	return &{TplTableNameCamelCase}Dao{M: d.M.Fields(fieldNamesOrMapStruct...)}
}

// FieldsEx 指定例外的字段，(不被操作的字段)，多个字段使用字符','连接。
// 参数<fieldNamesOrMapStruct>的类型可以是string/map/*map/struct/*struct。
func (d *{TplTableNameCamelCase}Dao) FieldsEx(fieldNamesOrMapStruct ...interface{}) *{TplTableNameCamelCase}Dao {
	return &{TplTableNameCamelCase}Dao{M: d.M.FieldsEx(fieldNamesOrMapStruct...)}
}

// Option 设置模型的“额外操作”选项。
func (d *{TplTableNameCamelCase}Dao) Option(option int) *{TplTableNameCamelCase}Dao {
	return &{TplTableNameCamelCase}Dao{M: d.M.Option(option)}
}

// OmitEmpty 空值过滤，(过滤输入参数中的空值: nil,"",0)。
func (d *{TplTableNameCamelCase}Dao) OmitEmpty() *{TplTableNameCamelCase}Dao {
	return &{TplTableNameCamelCase}Dao{M: d.M.OmitEmpty()}
}

// Filter 过滤提交参数中不符合表结构的数据项。
func (d *{TplTableNameCamelCase}Dao) Filter() *{TplTableNameCamelCase}Dao {
	return &{TplTableNameCamelCase}Dao{M: d.M.Filter()}
}

// Where 设置模型的条件语句。参数<where>可以是string/map/gmap/slice/struct/*struct等类型。
func (d *{TplTableNameCamelCase}Dao) Where(where interface{}, args ...interface{}) *{TplTableNameCamelCase}Dao {
	return &{TplTableNameCamelCase}Dao{M: d.M.Where(where, args...)}
}

// WherePri方法的功能同Where，但提供了对表主键的智能识别。
// 如果主键是“id”，并且给定<where>参数为“123”，则WherePri函数将条件视为“id=123”，而M.where将条件视为字符串“123”。
func (d *{TplTableNameCamelCase}Dao) WherePri(where interface{}, args ...interface{}) *{TplTableNameCamelCase}Dao {
	return &{TplTableNameCamelCase}Dao{M: d.M.WherePri(where, args...)}
}

// And 在where语句中添加“AND”条件。
func (d *{TplTableNameCamelCase}Dao) And(where interface{}, args ...interface{}) *{TplTableNameCamelCase}Dao {
	return &{TplTableNameCamelCase}Dao{M: d.M.And(where, args...)}
}

// Or 在where语句中添加“OR”条件。
func (d *{TplTableNameCamelCase}Dao) Or(where interface{}, args ...interface{}) *{TplTableNameCamelCase}Dao {
	return &{TplTableNameCamelCase}Dao{M: d.M.Or(where, args...)}
}

// Group 分组 (设置模型的“group by”语句)。
func (d *{TplTableNameCamelCase}Dao) Group(groupBy string) *{TplTableNameCamelCase}Dao {
	return &{TplTableNameCamelCase}Dao{M: d.M.Group(groupBy)}
}

// Order 排序 (设置模型的“order by”语句)。
func (d *{TplTableNameCamelCase}Dao) Order(orderBy ...string) *{TplTableNameCamelCase}Dao {
	return &{TplTableNameCamelCase}Dao{M: d.M.Order(orderBy...)}
}

// Limit 设置模型的“limit”语句。
// 参数<limit>可以是一个或两个数字，如果传递了两个数字，则为模型设置“limit limit[0]、limit[1]”语句，否则设置“limit limit[0]”语句。
func (d *{TplTableNameCamelCase}Dao) Limit(limit ...int) *{TplTableNameCamelCase}Dao {
	return &{TplTableNameCamelCase}Dao{M: d.M.Limit(limit...)}
}

// Offset 设置模型的“offset”语句。
// 它只适用于某些数据库，如SQLServer、PostgreSQL等。
func (d *{TplTableNameCamelCase}Dao) Offset(offset int) *{TplTableNameCamelCase}Dao {
	return &{TplTableNameCamelCase}Dao{M: d.M.Offset(offset)}
}

// Page 设置模型的页码。参数<page>从1开始分页。
// 注意，对于“Limit”语句，Limit函数从0开始是不同的。
func (d *{TplTableNameCamelCase}Dao) Page(page, limit int) *{TplTableNameCamelCase}Dao {
	return &{TplTableNameCamelCase}Dao{M: d.M.Page(page, limit)}
}

// Batch 指定批量操作中分批操作的条数数量 (默认是10)
func (d *{TplTableNameCamelCase}Dao) Batch(batch int) *{TplTableNameCamelCase}Dao {
	return &{TplTableNameCamelCase}Dao{M: d.M.Batch(batch)}
}

// Cache 设置模型的缓存功能，它缓存sql的结果。
func (d *{TplTableNameCamelCase}Dao) Cache(duration time.Duration, name ...string) *{TplTableNameCamelCase}Dao {
	return &{TplTableNameCamelCase}Dao{M: d.M.Cache(duration, name...)}
}

// Data 设置模型的操作数据，<data>可以是string/map/gmap/slice/struct/*struct等类型。
func (d *{TplTableNameCamelCase}Dao) Data(data ...interface{}) *{TplTableNameCamelCase}Dao {
	return &{TplTableNameCamelCase}Dao{M: d.M.Data(data...)}
}

// Take 通过M.WherePri和M.One检索并返回单个记录，并将结果返回为*model.{TplTableNameCamelCase}。
// 如果没有使用表中给定的条件检索到记录，则返回nil。另见M.WherePri和M.One。
func (d *{TplTableNameCamelCase}Dao) Take(where ...interface{}) (*model.{TplTableNameCamelCase}, error) {
	one, err := d.M.Take(where...)
	if err != nil {
		return nil, err
	}
	var entity *model.{TplTableNameCamelCase}
	if err = one.Struct(&entity); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return entity, nil
}

// Find 对模型执行“SELECT FROM…”语句，并将结果返回为[]*model.{TplTableNameCamelCase}
// 通过M.WherePri和M.All检索并返回结果集。另见M.WherePri和M.All。
func (d *{TplTableNameCamelCase}Dao) Find(where ...interface{}) ([]*model.{TplTableNameCamelCase}, error) {
	all, err := d.M.Find(where...)
	if err != nil {
		return nil, err
	}
	var entities []*model.{TplTableNameCamelCase}
	if err = all.Structs(&entities); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return entities, nil
}

// Scan 根据参数<pointer>的类型自动调用Struct或Structs函数
// <pointer>可以是*Struct/**Struct，也可以是*[]struct/*[]*struct。
func (d *{TplTableNameCamelCase}Dao) Scan(pointer interface{}, where ...interface{}) error {
	return d.M.Scan(pointer, where...)
}

// Chunk 使用给定的大小和回调函数迭代表。
func (d *{TplTableNameCamelCase}Dao) Chunk(limit int, callback func(entities []*model.{TplTableNameCamelCase}, err error) bool) {
	d.M.Chunk(limit, func(result orm.Result, err error) bool {
		var entities []*model.{TplTableNameCamelCase}
		err = result.Structs(&entities)
		if err == sql.ErrNoRows {
			return false
		}
		return callback(entities, err)
	})
}

// LockUpdate 用于创建FOR UPDATE锁，避免选择行被其它共享锁修改或删除，FOR UPDATE会阻塞其他锁定性读对锁定行的读取
func (d *{TplTableNameCamelCase}Dao) LockUpdate() *{TplTableNameCamelCase}Dao {
	return &{TplTableNameCamelCase}Dao{M: d.M.LockUpdate()}
}

// LockShared 使用LockShared方法，在运行sql语句时带一把”共享锁“，共享锁可以避免被选择的行被修改，直到事务提交。
func (d *{TplTableNameCamelCase}Dao) LockShared() *{TplTableNameCamelCase}Dao {
	return &{TplTableNameCamelCase}Dao{M: d.M.LockShared()}
}

// Delete 数据硬删除，被删除的数据不可恢复，请慎重使用。
func (d *{TplTableNameCamelCase}Dao) Unscoped() *{TplTableNameCamelCase}Dao {
	return &{TplTableNameCamelCase}Dao{M: d.M.Unscoped()}
}
`
