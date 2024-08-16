package daoctl

import (
	"context"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/kysion/base-library/utility/base_funs"
	"github.com/kysion/base-library/utility/daoctl/dao_interface"
)

// ContextModelTableKey 模型表名称的上下文键名。
const contextModelTableKey = "__TABLE__"

var (
	enableOrmCache = true
	extModelMap    = map[string]map[string]func(model *gdb.Model, conf *dao_interface.DaoConfig, data ...interface{}) *gdb.Model{}
)

// EnableOrmCache 启用所有的ORM缓存。
func EnableOrmCache() {
	enableOrmCache = true
}

// DisabledOrmCache 禁用所有的ORM缓存。
func DisabledOrmCache() {
	enableOrmCache = false
}

// MakeExtModelMap 用于注册外部模型的处理函数到指定的表。
// 该函数允许多个处理函数被同时注册到同一个表中。
// 参数 tableKey 指定需要注册的表名。
// 参数 f 为一个变长参数，代表一个或多个处理函数的指针。
// 每个处理函数接收一个模型对象、一个配置对象以及任意数量的额外数据，并返回一个处理后的模型对象。
func MakeExtModelMap(tableKey string, f ...map[string]func(model *gdb.Model, conf *dao_interface.DaoConfig, data ...interface{}) *gdb.Model) {
	// 遍历所有传入的处理函数指针
	for _, v := range f {
		// 将每个处理函数指针注册到内部的扩展模型映射中
		extModelMap[tableKey] = v
	}
}

// NewDaoConfig 创建并初始化一个DaoConfig实例。
// 参数:
// - ctx: 上下文对象，用于传递请求范围的数据。
// - dao: 实现了IDao接口的对象。
// - cacheOption: 可选的缓存配置。
// 返回值:
// - 初始化后的DaoConfig实例。
func NewDaoConfig(ctx context.Context, dao dao_interface.IDao, cacheOption ...*gdb.CacheOption) dao_interface.DaoConfig {
	result := dao_interface.DaoConfig{
		Dao:   dao,
		DB:    dao.DB(),
		Table: dao.Table(),
		Group: dao.Group(),
	}

	// 根据数据访问对象（DAO）的表名，检查是否存在扩展查询条件。
	// 如果扩展查询条件映射不为空，并且针对当前DAO表的扩展查询条件存在，
	// 则将这些扩展查询条件应用到结果对象上。
	if extModelMap != nil && extModelMap[dao.Table()] != nil {
		if result.ExtWhere == nil {
			result.ExtWhere = make(map[string]func(model *gdb.Model, conf *dao_interface.DaoConfig, data ...interface{}) *gdb.Model, 0)
		}
		for k, item := range extModelMap[dao.Table()] {
			result.ExtWhere[k] = item
		}
		//_ = gconv.MapToMap(dao_interface.ExtModelMap[dao.Table()], &result.ExtWhere)
	}

	// 设置上下文中的表名。
	ctx = context.WithValue(ctx, contextModelTableKey, result.Table)
	// 设置上下文中特定表的配置。
	ctx = context.WithValue(ctx, result.Table, &result)

	// 根据上下文和表名初始化数据库模型。
	result.Model = dao.DB().Model(dao.Table()).Safe().Ctx(ctx)

	// 获取配置中指定的忽略缓存的表列表。
	dataCacheConf := g.Cfg().MustGet(ctx, "ormCache.ignore.tables")
	// 如果启用了缓存，并且配置中未指定忽略缓存的表，则检查当前表是否被忽略。
	if enableOrmCache == true && dataCacheConf.String() != "*" && !IsIgnoreOrmCacheByCtx(ctx, dao.Table()) {
		// 从配置中获取忽略缓存的表列表。
		cacheIgnoreTables := g.Cfg().MustGet(ctx, "ormCache.ignore.tables").Strings()
		// 如果存在忽略缓存的表列表，则过滤空字符串和重复的条目。
		if len(cacheIgnoreTables) > 0 {
			// 去除重复的条目。
			cacheIgnoreTables = base_funs.Unique(cacheIgnoreTables)
			// 过滤空字符串和重复的条目。
			cacheIgnoreTables = base_funs.FilterEmpty(cacheIgnoreTables)
		}
		// 如果当前表不在忽略缓存列表中，则配置缓存选项。
		if result.IsIgnoreCache() == false || !base_funs.Contains(cacheIgnoreTables, dao.Table()) {
			if len(cacheOption) == 0 {
				// 如果没有提供缓存选项，则自动生成一个。
				result.CacheOption = MakeDaoCache(dao.Table())
				result.Model = result.Model.Cache(*result.CacheOption)
			} else if cacheOption[0] != nil {
				// 如果提供了缓存选项，则使用提供的选项。
				result.CacheOption = cacheOption[0]
				result.Model = result.Model.Cache(*result.CacheOption)
			}
			// 注册DAO钩子，以便在数据库操作前后执行自定义逻辑。
			result.Model = RegisterDaoHook(result.Model)
		}
	}

	return result
}

// ExecExWhere 执行扩展的条件查询。
// 该函数会根据模型的上下文信息，动态地对模型应用一系列扩展的查询条件。
// 参数 model 是要进行查询的模型。
// 参数 data 是传递给扩展查询函数的额外数据。
// 返回值是应用了扩展查询条件后的模型。
func ExecExWhere(model *gdb.Model, data ...interface{}) *gdb.Model {
	// 从模型的上下文中获取表名。
	if tableName, ok := model.GetCtx().Value(contextModelTableKey).(string); ok {
		// 根据表名生成上下文中的扩展条件查询键名。
		ctxWhereKey := "_ctx_ext_where_key_" + tableName

		// 尝试从上下文中获取已存在的扩展条件查询键名列表。
		existingKeys, ok := model.GetCtx().Value(ctxWhereKey).([]string)

		// 如果获取失败，将 existingKeys 设置为 nil。
		if !ok {
			existingKeys = nil
		}

		// 从上下文中获取与表名相关联的 DaoConfig 对象。
		if conf, ok := model.GetCtx().Value(tableName).(*dao_interface.DaoConfig); ok && conf != nil {
			// 如果 DaoConfig 中定义了扩展的条件查询函数。
			if conf.ExtWhere != nil {
				// 遍历扩展的条件查询函数列表，并依次调用。
				for k, v := range conf.ExtWhere {
					// 检查当前扩展条件查询函数的键名是否已存在于 existingKeys 中，如果是，则跳过。
					if existingKeys != nil && base_funs.Contains(existingKeys, k) {
						continue
					}

					// 调用每个扩展条件查询函数，并更新 model。
					// 注意：这里假设 v 是一个函数，其签名为 func(*gdb.Model, *dao_interface.DaoConfig, ...interface{}) *gdb.Model。
					model = v(model, conf, data...)
				}
			}
		}
	}

	// 返回执行完扩展条件查询后的模型。
	return model
}

// IgnoreExtModel 该函数用于在上下文中忽略指定的扩展模型字段。
// 它接受一个上下文对象和一个表名字符串，以及一个可变长的字符串数组作为条件键。
// 如果没有指定条件键，它将直接返回原始上下文。
// 如果指定了条件键，它会将这些键添加到上下文中，以便在后续处理中可以忽略这些字段。
// 参数:
//
//	ctx - 上下文对象，用于存储请求相关的数据。
//	tableName - 表名字符串，用于区分不同的数据表。
//	whereKey... - 可变长的字符串数组，代表要忽略的条件键。
//
// 返回值:
//
//	返回更新后的上下文对象，其中包含了要忽略的条件键列表。
func IgnoreExtModel(ctx context.Context, tableName string, whereKey ...string) context.Context {
	// 检查是否提供了条件键，如果没有，则直接返回原始上下文
	if len(whereKey) == 0 {
		return ctx
	}

	// 定义一个上下文键，用于存储当前表需要忽略的条件键列表
	ctxWhereKey := "_ctx_ext_where_key_" + tableName

	// 从上下文中获取现有键列表，如果不存在则初始化为空切片
	existingKeys, ok := ctx.Value(ctxWhereKey).([]string)
	if !ok || existingKeys == nil {
		existingKeys = make([]string, 0)
	}

	// 遍历传入的条件键，添加到现有键列表中，同时避免重复键
	for _, key := range whereKey {
		if !base_funs.Contains(existingKeys, key) {
			existingKeys = append(existingKeys, key)
		}
	}

	// 移除键列表中的空字符串，保持列表的清洁
	existingKeys = base_funs.FilterEmpty(existingKeys)

	// 将更新后的键列表放回上下文中，并返回更新后的上下文对象
	return context.WithValue(ctx, ctxWhereKey, existingKeys)
}

// IsIgnoreExtModel 检查给定的表名和条件键是否在上下文中被标记为需要忽略。
//
// 参数:
//
//	ctx - 上下文对象，用于传递忽略条件键的相关信息。
//	tableName - 需要检查的表名。
//	whereKey - 需要检查的条件键。
//
// 返回值:
//
//	如果表名和条件键组合被标记为需要忽略，则返回true；否则返回false。
func IsIgnoreExtModel(ctx context.Context, tableName string, whereKey string) bool {
	// 定义一个上下文键，用于存储当前表需要忽略的条件键列表
	ctxWhereKey := "_ctx_ext_where_key_" + tableName

	// 从上下文中获取现有键列表
	existingKeys, ok := ctx.Value(ctxWhereKey).([]string)
	if !ok || existingKeys == nil {
		return false
	}

	// 检查给定的条件键是否存在于现有键列表中
	return base_funs.Contains(existingKeys, whereKey)
}

const contextOrmTableCacheKey = "_ctx_orm_cache_key_"

// IsIgnoreOrmCacheByCtx 检查给定的表名和操作类型是否在上下文中被标记为需要忽略缓存。
//
// 参数:
//
//	ctx - 上下文对象，用于传递忽略缓存操作的相关信息。
//	tableName - 需要检查的表名。
//	action - 表的操作类型。
//
// 返回值:
//
//	如果表名和操作类型组合被标记为需要忽略缓存，则返回true；否则返回false。
func IsIgnoreOrmCacheByCtx(ctx context.Context, tableName string) bool {
	// 从现有context中获取存储表名的sync.Map
	existingTableNames, ok := ctx.Value(contextOrmTableCacheKey).([]string)
	if !ok || existingTableNames == nil {
		return false
	}

	// 如果找到匹配项，则表示该操作被标记为忽略缓存
	return base_funs.Contains(existingTableNames, tableName)
}

// SetTableIgnoreOrmCacheByCtx 根据指定的表名列表，在上下文中记录需要忽略ORM缓存操作的信息。
// 参数:
//
//	ctx: 传入的上下文对象，用于存储和传递忽略ORM缓存操作的信息。
//	tableNames: 可变参数，指定需要忽略ORM缓存操作的表名列表。
//
// 返回值:
//
//	返回一个新的上下文对象，其中包含了更新后的忽略ORM缓存操作的信息。
func SetTableIgnoreOrmCacheByCtx(ctx context.Context, tableNames ...string) context.Context {
	// 如果没有指定任何表操作，则直接返回原始context
	if len(tableNames) == 0 {
		return ctx
	}

	// 从现有context中加载或初始化存储表名的slice
	existingTableNames, ok := ctx.Value(contextOrmTableCacheKey).([]string)
	if !ok || existingTableNames == nil {
		existingTableNames = make([]string, 0)
	}

	// 合并传入的表名到现有的表名列表中，并确保列表中没有重复的表名
	newExistingTableNames := append(existingTableNames, tableNames...)
	// 去除重复的表名
	newExistingTableNames = base_funs.Unique(newExistingTableNames)
	// 移除空字符串，保持列表的清洁
	newExistingTableNames = base_funs.FilterEmpty(newExistingTableNames)

	// 返回包含忽略缓存数据相关Dao操作信息的新context
	return context.WithValue(ctx, contextOrmTableCacheKey, newExistingTableNames)
}
