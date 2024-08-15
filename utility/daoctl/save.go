package daoctl

import (
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/kysion/base-library/utility/daoctl/dao_interface"
)

// Save 保存模型数据。
//
// model 参数是指向要保存的模型的指针，data 是与模型相关的数据。
// 该函数首先执行扩展的条件查询，然后尝试保存模型数据。
// 如果保存操作成功，它将返回受影响的行数；如果失败，则返回 0。
//
// 参数:
// - model: 要保存的模型指针。
// - data: 与模型相关的数据，用于更新或插入操作。
//
// 返回值:
// - rowsAffected: 受影响的行数，成功时为实际受影响的行数，失败时为 0。
func Save(model *gdb.Model, data ...interface{}) (rowsAffected int64) {
	// 执行扩展的条件查询，以确保数据满足保存的条件。
	model = dao_interface.ExecExWhere(model, data...)

	// 尝试保存模型数据。
	result, err := model.Save(data...)

	// 如果保存操作出错，返回 0 表示没有受影响的行。
	if err != nil {
		return 0
	}

	// 获取并返回保存操作受影响的行数。
	rowsAffected, err = result.RowsAffected()
	return rowsAffected
}

// SaveWithError 将数据保存到模型中，并返回影响的行数和可能的错误。
// 该函数使用 ExecExWhere 对模型执行额外的处理，然后尝试保存数据。
// 如果保存过程中出现错误，将返回错误本身；如果没有错误，将返回受影响的行数。
//
// 参数:
// - model: 待操作的gdb.Model对象指针。
// - data: 可选的附加数据，用于补充模型信息或操作条件。
//
// 返回值:
// - rowsAffected: 保存操作影响的行数。
// - err: 保存操作中可能出现的错误，如果没有错误，则为nil。
func SaveWithError(model *gdb.Model, data ...interface{}) (rowsAffected int64, err error) {
	// 对模型执行额外的处理，这是为了确保数据在保存前满足特定条件或逻辑。
	model = dao_interface.ExecExWhere(model, data...)

	// 尝试保存模型数据。
	result, err := model.Save(data...)

	// 如果保存过程中出现错误，返回错误。
	if err != nil {
		return 0, err
	}

	// 返回保存操作影响的行数。
	return result.RowsAffected()
}
