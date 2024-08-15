package daoctl

import (
	"github.com/gogf/gf/v2/database/gdb"
)

// Scan 从数据库中查询数据，并将其反序列化为指定类型的实例。
// 它接受一个 *gdb.Model 对象作为查询模板，并返回查询结果。
// 如果查询成功，它返回结果实例；如果查询失败，它返回 nil。
// 参数:
//
//	model: *gdb.Model 类型的查询模板。
//
// 返回值:
//
//	*T: 查询结果的指针，如果查询失败则为 nil。
func Scan[T any](model *gdb.Model) *T {
	// 执行带有额外条件的查询，这些条件可能在查询执行前需要设置。
	model = ExecExWhere(model)

	// 创建一个 T 类型的空实例，用于存储查询结果。
	result := new(T)
	// 使用模型的 Scan 方法从数据库中加载数据到 result 实例中。
	if err := model.Scan(result); err != nil {
		// 如果 Scan 方法返回错误，表明查询失败，此时返回 nil。
		return nil
	}
	// 如果查询成功，返回存储结果的实例。
	return result
}

// ScanWithError 用于从数据库查询结果中扫描并返回一个特定类型的实例，同时处理可能发生的错误。
// 这个函数接收一个gdb.Model指针作为输入参数，该模型已经设置好查询条件。
// 函数返回的是一个泛型T的指针和一个错误类型。
// T是任何符合泛型约束的类型，这里使用了泛型编程来实现类型通用性。
// 参数:
//
//	model (*gdb.Model): 用于从数据库中查询数据的模型，包含查询条件。
//
// 返回值:
//
//	*T: 查询结果转换为的具体类型实例，如果查询失败或者转换失败，返回nil。
//	error: 如果查询或者转换过程中出现错误，返回具体的错误信息，否则返回nil。
func ScanWithError[T any](model *gdb.Model) (*T, error) {
	// ExecExWhere 是一个假设存在的函数，用于执行额外的查询优化或者条件添加操作。
	// 这里不详细展开其具体实现，因为它不是这段代码关注的焦点。
	model = ExecExWhere(model)

	// new(T) 用于创建一个T类型的零值实例，用于接下来接收查询结果。
	result := new(T)
	// 使用model的Scan方法将查询结果填充到result中，如果出现错误则返回错误。
	if err := model.Scan(result); err != nil {
		return nil, err
	}
	// 如果一切顺利，返回填充好的result实例和nil错误。
	return result, nil
}

// ScanList 根据给定的模型和属性名，扫描列表数据并绑定到指定的结构体。
// 该函数支持泛型，可以处理多种类型的列表数据。
// 参数:
//
//	model: *gdb.Model 类型，表示数据库模型。
//	bindToAttrName: 字符串类型，表示要绑定数据的结构体属性名。
//	relationAttrNameAndFields: 可变参数，用于指定关联属性名和其他字段。
//
// 返回值:
//
//	返回绑定数据后的结构体指针，如果出现错误则返回 nil。
func ScanList[T any](model *gdb.Model, bindToAttrName string, relationAttrNameAndFields ...string) *T {
	// ExecExWhere 方法用于执行额外的查询条件，确保数据的准确性。
	model = ExecExWhere(model)

	// 初始化结果变量，使用泛型 T 的零值。
	result := new(T)
	// 使用 ScanList 方法从模型中扫描数据，并绑定到 result 上。
	// 如果扫描过程中出现错误，则返回 nil。
	if err := model.ScanList(result, bindToAttrName, relationAttrNameAndFields...); err != nil {
		return nil
	}
	// 返回成功绑定数据后的结果。
	return result
}

// ScanListWithError 用于处理模型查询结果，并将其扫描到指定的结构体中。
// 此函数提供了一个泛型实现，可以处理多种类型的模型。
// 参数：
//
//	model: 指向数据库模型的指针，用于执行查询操作。
//	bindToAttrName: 字符串，指定将结果绑定到的属性名称。
//	relationAttrNameAndFields: 可选参数，用于指定关联属性名称及其字段。
//
// 返回值：
//
//	*T: 查询结果，被扫描到一个新实例中。
//	error: 如果在执行查询或扫描过程中发生错误，将返回该错误。
func ScanListWithError[T any](model *gdb.Model, bindToAttrName string, relationAttrNameAndFields ...string) (*T, error) {
	// 对模型执行额外的 where 条件
	model = ExecExWhere(model)

	// 初始化结果变量，使用泛型 T
	result := new(T)
	// 执行扫描操作，将查询结果填充到 result 中
	if err := model.ScanList(result, bindToAttrName, relationAttrNameAndFields...); err != nil {
		// 如果扫描过程中出现错误，返回错误信息
		return nil, err
	}
	// 返回成功扫描的结果和空错误
	return result, nil
}
