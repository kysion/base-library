package daoctl

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/kysion/base-library/base_consts"
	"time"
	"unsafe"
)

// 定义一个钩子处理程序，用于处理不同类型的数据库操作。
// 该处理程序通过清洁缓存来响应更新、插入和删除操作，
// 并通过查询输入的Next方法来响应选择操作。
var HookHandler = gdb.HookHandler{
	// 使用cleanCache函数来处理更新操作
	Update: cleanCache[gdb.HookUpdateInput],
	// 使用cleanCache函数来处理插入操作
	Insert: cleanCache[gdb.HookInsertInput],
	// 使用cleanCache函数来处理删除操作
	Delete: cleanCache[gdb.HookDeleteInput],
	// 定义选择操作的处理逻辑
	Select: func(ctx context.Context, in *gdb.HookSelectInput) (result gdb.Result, err error) {
		// 调用输入的Next方法执行查询，并返回结果和可能的错误
		result, err = in.Next(ctx)
		return
	},
}

// iHookInput 是一个接口，定义了钩子输入的结构。
// 它主要用于判断操作是否属于事务，并提供执行下一个操作的能力。
type iHookInput interface {
	// IsTransaction 返回当前操作是否属于事务的标志。
	IsTransaction() bool

	// Next 执行下一个钩子操作。
	// 参数：
	// - ctx: 上下文，用于传递请求的取消信号、超时设置等。
	// 返回值：
	// - result: 操作结果，具体类型依赖于操作。
	// - err: 操作中出现的错误，如果没有错误则为nil。
	Next(ctx context.Context) (result sql.Result, err error)
}

// cleanCache 清理缓存函数，根据不同的输入类型（插入、更新、删除）来清理相应的数据库缓存。
// 输入参数 T 可以是 gdb.HookInsertInput、gdb.HookUpdateInput 或 gdb.HookDeleteInput 类型。
// 返回值 result 为清理缓存后的数据库操作结果，err 为可能出现的错误。
func cleanCache[T gdb.HookInsertInput | gdb.HookUpdateInput | gdb.HookDeleteInput](ctx context.Context, in *T) (result sql.Result, err error) {
	// 将输入参数 in 转换为 iHookInput 接口类型，以便调用 Next 方法。
	v, ok := interface{}(in).(iHookInput)
	if !ok {
		// 如果转换失败，返回错误。
		return result, fmt.Errorf("input does not implement iHookInput")
	}

	var table string
	var model *gdb.Model

	// 定义缓存清理配置，-1 表示使用全局配置，false 表示不强制清理。
	conf := gdb.CacheOption{
		Duration: -1,
		Force:    false,
	}

	// 根据输入参数的不同类型，清理相应的缓存。
	if input, ok := interface{}(in).(*gdb.HookInsertInput); ok == true {
		input.Model.Cache(conf)
		table = input.Table
		model = input.Model
	} else if input, ok := interface{}(in).(*gdb.HookUpdateInput); ok == true {
		input.Model.Cache(conf)
		table = input.Table
		model = input.Model
	} else if input, ok := interface{}(in).(*gdb.HookDeleteInput); ok == true {
		input.Model.Cache(conf)
		table = input.Table
		model = input.Model
	}

	// 清理完缓存后，如果有表名，则进行表名的格式化处理。
	if table != "" {
		table = gstr.SplitAndTrim(table, " ")[0]
		table = gstr.SplitAndTrim(table, ",")[0]
		table = gstr.Replace(table, "\"", "")
	}

	// 根据表名从缓存中移除对应的缓存项。
	if model != nil {
		db := *(*gdb.DB)(unsafe.Pointer(model))

		cacheKeys, _ := db.GetCache().KeyStrings(ctx)
		for _, key := range cacheKeys {
			if gstr.HasPrefix(key, table) || gstr.HasPrefix(key, "SelectCache:"+table) || gstr.HasPrefix(key, "SelectCache:default@#"+table) {
				db.GetCache().Remove(db.GetCtx(), key)
			}
		}
	}

	// 执行下一步操作并返回结果和错误。
	result, err = v.Next(ctx)
	return
}

// RemoveQueryCache 用于根据前缀移除查询缓存。
// 该函数获取所有缓存的键，然后遍历这些键，如果键以指定的前缀开始，
// 或者以"SelectCache:"加上前缀开始，或者以"SelectCache:default@#"加上前缀开始，
// 则移除该缓存项。这主要用于在数据源发生变化时，清理相关的缓存，
// 以确保数据的一致性和新鲜度。
//
// 参数:
// - db: 数据库对象，用于访问缓存。
// - prefix: 缓存键的前缀，用于匹配和移除缓存项。
func RemoveQueryCache(db gdb.DB, prefix string) {
	// 获取所有缓存的键
	cacheKeys, _ := db.GetCache().KeyStrings(db.GetCtx())
	// 遍历缓存键
	for _, key := range cacheKeys {
		// 判断缓存键是否匹配移除条件
		if gstr.HasPrefix(key, prefix) || gstr.HasPrefix(key, "SelectCache:"+prefix) || gstr.HasPrefix(key, "SelectCache:default@#"+prefix) {
			// 移除匹配的缓存项
			db.GetCache().Remove(db.GetCtx(), key)
		}
	}
}

// MakeDaoCache 根据全局配置和指定表名创建数据库缓存配置对象。
// 该函数通过查找全局缓存配置列表，找到与给定表名匹配的配置，
// 并据此配置缓存的过期时间和是否强制使用缓存。
// 参数table：要设置缓存的表名。
// 返回值：*gdb.CacheOption 类型的缓存配置对象，用于数据库操作中启用缓存。
func MakeDaoCache(table string) *gdb.CacheOption {
	// 初始化默认的缓存配置，缓存默认过期时间为24小时，且不强制使用缓存。
	conf := &gdb.CacheOption{
		Duration: time.Hour * 24,
		Force:    false,
	}

	// 遍历全局缓存配置列表，寻找与参数table匹配的配置。
	for _, cacheConf := range base_consts.Global.OrmCacheConf {
		// 当找到匹配的表名时，更新缓存的过期时间和强制使用缓存的配置。
		if cacheConf.TableName == table {
			// 根据配置设置缓存过期时间。
			conf.Duration = time.Second * (time.Duration)(cacheConf.ExpireSeconds)
			// 根据配置设置是否强制使用缓存。
			conf.Force = cacheConf.Force
		}
	}

	// 返回配置好的缓存选项。
	return conf
}

// RegisterDaoHook 注册数据库操作前的钩子函数。
// 该函数将 HookHandler 注册为数据库模型的操作钩子，以便在数据操作执行前后执行特定的逻辑。
// 参数:
//
//	model: 指向数据库模型的指针。
//
// 返回值:
//
//	返回注册了钩子函数后的模型指针。
func RegisterDaoHook(model *gdb.Model) *gdb.Model {
	return model.Hook(HookHandler)
}
