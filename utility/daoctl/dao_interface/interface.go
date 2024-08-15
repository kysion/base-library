package dao_interface

import (
	"context"
	"github.com/gogf/gf/v2/database/gdb"
)

// ContextModelTableKey 模型表名称的上下文键名。
const ContextModelTableKey = "__TABLE__"

var ExtModelMap = map[string]map[string]func(model *gdb.Model, conf *DaoConfig, data ...interface{}) *gdb.Model{}

// DaoConfig 定义了数据访问对象（DAO）的配置结构。
type DaoConfig struct {
	Dao         IDao                                                                               // IDao接口实例，用于数据库操作。
	DB          gdb.DB                                                                             // 数据库连接实例。
	Table       string                                                                             // 数据库操作的目标表名。
	Group       string                                                                             // 数据库操作的分组名。
	Model       *gdb.Model                                                                         // 数据库模型，用于执行数据库查询。
	CacheOption *gdb.CacheOption                                                                   // 缓存选项，用于配置查询缓存。
	HookHandler *gdb.HookHandler                                                                   // 钩子处理器，用于在数据库操作前后执行自定义逻辑。
	ignoreCache bool                                                                               // 标志位，指示是否忽略缓存。
	ExtWhere    map[string]func(model *gdb.Model, conf *DaoConfig, data ...interface{}) *gdb.Model // 在查询之前应用的额外查询条件。
}

// IDao 定义了数据访问对象（DAO）的基本接口。
type IDao interface {
	DB() gdb.DB                                                                                // 获取数据库连接实例。
	Table() string                                                                             // 获取操作的表名。
	Group() string                                                                             // 获取数据库操作的分组名。
	Ctx(ctx context.Context, cacheOption ...*gdb.CacheOption) *gdb.Model                       // 在给定上下文中创建模型，并可选配置缓存选项。
	Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) // 在给定上下文中执行事务操作。
	DaoConfig(ctx context.Context, cacheOption ...*gdb.CacheOption) *DaoConfig                 // 获取DAO的配置。
	GetExtWhereKeys() []string
	IsIgnoreCache() bool
}

// TIDao 是泛型接口，扩展了IDao接口，增加了Columns方法以支持泛型操作。
type TIDao[TColumns any] interface {
	IDao               // 继承IDao接口。
	Columns() TColumns // 获取泛型类型的列信息。
	IgnoreCache() IDao
	IgnoreExtModel(whereKey ...string) IDao
}

// IgnoreCache 方法用于标记是否忽略缓存。
// 返回值:
// - 修改后的DaoConfig实例。
func (d *DaoConfig) IgnoreCache() *DaoConfig {
	d.ignoreCache = true
	return d
}

func (d *DaoConfig) IsIgnoreCache() bool {
	return d.ignoreCache
}

// IgnoreExtModel 方法用于忽略在查询之前应用的指定额外查询条件。
// 参数:
// - whereKey: 要忽略的查询条件键名列表。
// 返回值:
// - 修改后的DaoConfig实例。
func (d *DaoConfig) IgnoreExtModel(whereKey ...string) *DaoConfig {
	if d.ExtWhere != nil {
		for _, key := range whereKey {
			delete(d.ExtWhere, key) // 删除指定的查询条件。
		}
	}
	return d
}
