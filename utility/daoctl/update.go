package daoctl

import (
	"github.com/gogf/gf/v2/database/gdb"
)

// Update 根据给定的条件更新模型，并返回受影响的行数。
// 该函数首先使用ExecExWhere函数执行带有条件的更新操作，然后调用model的Update方法进行更新操作。
// 参数model是待更新的模型，dataAndWhere是用于更新的数据和where条件。
// 返回值rowsAffected是更新操作影响的行数。
func Update(model *gdb.Model, dataAndWhere ...interface{}) (rowsAffected int64) {
	// 使用ExecExWhere函数执行带有条件的更新操作。
	model = ExecExWhere(model, dataAndWhere...)

	// 调用model的Update方法进行更新操作，返回更新结果和可能的错误。
	result, err := model.Update(dataAndWhere...)

	// 如果更新操作出现错误，则返回0表示没有行受到影响。
	if err != nil {
		return 0
	}

	// 获取更新操作影响的行数，并返回。
	rowsAffected, err = result.RowsAffected()
	return rowsAffected
}

// UpdateWithError 更新模型数据，并返回影响的行数和可能的错误。
// 该函数专注于处理带有额外执行条件的更新操作，提供了一种灵活的数据更新方式。
// 参数:
//
//	model: 待更新的模型指针。
//	dataAndWhere: 用于更新的数据和作为更新条件的where子句，可变参数。
//
// 返回值:
//
//	rowsAffected: 更新操作影响的行数。
//	err: 更新操作过程中可能产生的错误。
func UpdateWithError(model *gdb.Model, dataAndWhere ...interface{}) (rowsAffected int64, err error) {
	// 执行可能的额外条件（如软删除等）。
	model = ExecExWhere(model)

	// 执行更新操作。
	result, err := model.Update(dataAndWhere...)
	if err != nil {
		// 如果更新过程中出现错误，返回错误信息。
		return 0, err
	}
	// 返回影响的行数。
	return result.RowsAffected()
}
