package dao_interface

import (
	"github.com/gogf/gf/v2/database/gdb"
)

// MakeExtModelMap 用于注册外部模型的处理函数到指定的表。
// 该函数允许多个处理函数被同时注册到同一个表中。
// 参数 tableKey 指定需要注册的表名。
// 参数 f 为一个变长参数，代表一个或多个处理函数的指针。
// 每个处理函数接收一个模型对象、一个配置对象以及任意数量的额外数据，并返回一个处理后的模型对象。
func MakeExtModelMap(tableKey string, f ...map[string]func(model *gdb.Model, conf *DaoConfig, data ...interface{}) *gdb.Model) {
	// 遍历所有传入的处理函数指针
	for _, v := range f {
		// 将每个处理函数指针注册到内部的扩展模型映射中
		ExtModelMap[tableKey] = v
	}
}

// ExecExWhere 执行扩展的条件查询。
//
// 该函数接收一个 *gdb.Model 类型的 model 参数和一个可变参数 data。其主要功能是在模型上执行扩展的条件查询。
// 如果存在扩展的条件查询函数（ExtWhere），则会依次调用它们，并将模型和数据传递给这些函数。
//
// 参数:
// - model: 一个指向 gdb.Model 的指针，代表数据库的模型。
// - data: 一个可变参数，包含传递给扩展条件查询函数的数据。
//
// 返回值:
// - *gdb.Model: 返回执行了扩展条件查询后的模型。
func ExecExWhere(model *gdb.Model, data ...interface{}) *gdb.Model {
	// 从模型的上下文中获取表名。
	if tableName, ok := model.GetCtx().Value(ContextModelTableKey).(string); ok {
		// 从上下文中获取与表名相关联的 DaoConfig 对象。
		if conf, ok := model.GetCtx().Value(tableName).(*DaoConfig); ok && conf != nil {
			// 如果 DaoConfig 中定义了扩展的条件查询函数。
			if conf.ExtWhere != nil {
				// 遍历扩展的条件查询函数列表，并依次调用。
				for _, v := range conf.ExtWhere {
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
