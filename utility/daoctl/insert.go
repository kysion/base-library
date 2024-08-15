package daoctl

import (
	"github.com/gogf/gf/v2/database/gdb"
)

// Insert 执行插入操作，并返回影响的行数。
// 它首先通过ExecExWhere对model进行额外的处理，然后执行插入操作。
// 参数model指向待操作的数据库模型，data为可变参数，用于ExecExWhere函数的额外处理和插入操作。
// 返回值rowsAffected表示插入操作影响的行数。
func Insert(model *gdb.Model, data ...interface{}) (rowsAffected int64) {
	// 对model应用ExecExWhere处理，以便进行额外的条件筛选或修改。
	model = ExecExWhere(model, data...)

	// 执行插入操作，并捕获可能的错误。
	result, err := model.Insert(data...)

	// 如果插入操作出现错误，则返回0，表示没有行受到影响。
	if err != nil {
		return 0
	}

	// 获取插入操作影响的行数，并捕获可能的错误。
	rowsAffected, err = result.RowsAffected()

	// 返回影响的行数。
	return rowsAffected
}

// InsertWithError 执行带有错误处理的插入操作。
// 它首先执行一个带有条件的删除操作，然后尝试插入给定的数据。
// 如果插入操作出错，它会返回错误。
//
// 参数:
// - model: 要操作的模型。
// - data: 待插入的数据，可以是多个接口{}类型的参数。
//
// 返回值:
// - rowsAffected: 受影响的行数。
// - err: 插入操作过程中可能发生的错误。
func InsertWithError(model *gdb.Model, data ...interface{}) (rowsAffected int64, err error) {
	// 使用传入的data参数执行一个带有条件的删除操作。
	// 这一步是为了确保在插入之前，根据提供的数据条件删除可能存在的旧数据。
	model = ExecExWhere(model, data...)

	// 尝试使用model插入data参数表示的数据。
	// 这一步是实际的数据插入操作，如果数据格式或数据库约束条件不满足，可能会产生错误。
	result, err := model.Insert(data...)

	// 如果插入操作出错，返回错误。
	// 通过这种方式，错误可以被上层调用者捕获和处理。
	if err != nil {
		return 0, err
	}

	// 返回插入操作影响的行数。
	// 这个数字可用于判断操作是否有效，以及有多少行被插入。
	return result.RowsAffected()
}

// InsertIgnore 执行一个“插入或忽略”操作，这意味著如果插入的数据已经存在，则不会执行插入。
// 这个函数使用 ExecExWhere 方法执行额外的 WHERE 条件筛选后，再执行插入操作，以确保数据的唯一性。
// 参数:
//
//	model: 指向要操作的模型。
//	data: 一个或多个要插入的数据项。
//
// 返回值:
//
//	rowsAffected: 受影响的行数，即成功插入的新数据项的数量。
func InsertIgnore(model *gdb.Model, data ...interface{}) (rowsAffected int64) {
	// 使用 ExecExWhere 方法对 model 进行额外的 WHERE 条件筛选，以确保数据的唯一性。
	model = ExecExWhere(model, data...)

	// 使用 InsertIgnore 方法尝试插入数据，这会自动忽略已存在的数据。
	result, err := model.InsertIgnore(data...)

	// 如果插入过程中出现错误，则返回受影响的行数为 0。
	if err != nil {
		return 0
	}

	// 获取并返回成功插入的行数。
	rowsAffected, err = result.RowsAffected()
	return rowsAffected
}

// InsertIgnoreWithError 执行有条件插入操作。
// 该函数尝试将数据插入到模型所表示的数据库表中，如果插入失败（例如，由于唯一约束冲突），则返回受影响的行数和错误。
// 使用 ExecExWhere 方法先执行一个排除条件的查询，然后使用结果作为插入操作的一部分。
// 参数 model 是待操作的模型，data 是待插入的数据，可以是多个。
// 返回值 rowsAffected 表示插入操作受影响的行数，err 表示插入操作过程中可能发生的错误。
func InsertIgnoreWithError(model *gdb.Model, data ...interface{}) (rowsAffected int64, err error) {
	// ExecExWhere 方法用于执行一个排除条件的查询，这里是为了在插入之前进行一些预处理操作。
	model = ExecExWhere(model, data...)

	// 尝试执行插入操作，这里使用的是 Model.Insert 方法，它允许插入多条记录。
	result, err := model.Insert(data...)

	// 如果插入操作发生错误，直接返回错误信息。
	if err != nil {
		return 0, err
	}

	// 插入成功后，返回受影响的行数。
	return result.RowsAffected()
}
