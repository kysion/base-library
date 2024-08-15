package daoctl

import (
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
)

// Delete 删除数据库中的记录。
//
// 参数:
// - model: 指向要删除记录的模型对象，该对象包含了删除操作所需的表信息和条件。
//
// 返回值:
// - rowsAffected: 删除操作影响的行数，是一个 int64 类型的值。
// - err: 删除操作过程中可能遇到的错误，如果没有错误发生，则返回 nil。
//
// 该函数首先通过调用 ExecExWhere 方法对模型进行额外的处理，这可能包括添加额外的
// 删除条件。然后，它尝试删除处理后的模型所代表的数据库记录。如果删除操作成功，
// 它将返回受影响的行数，否则将返回遇到的错误。
func Delete(model *gdb.Model) (rowsAffected int64, err error) {
	// 对模型执行额外的处理，可能是添加或修改删除条件。
	model = ExecExWhere(model)

	// 尝试根据模型删除数据库中的记录。
	result, err := model.Delete()
	if err != nil {
		return 0, err
	}

	// 获取删除操作影响的行数。
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return 0, err
	}

	// 返回受影响的行数和可能的错误。
	return rowsAffected, nil
}

// DeleteWithError 删除给定的模型，并返回受影响的行数和可能的错误。
// 参数 model: 指向要删除的模型的指针。
// 返回值:
// - rowsAffected: 受影响的行数，即被删除的行数。
// - err: 删除操作中可能出现的错误。
func DeleteWithError(model *gdb.Model) (rowsAffected int64, err error) {
	// 空指针检查
	if model == nil {
		return 0, fmt.Errorf("model is nil")
	}

	// 执行 ExecExWhere 并检查返回值
	updatedModel := ExecExWhere(model)
	if updatedModel == nil {
		return 0, fmt.Errorf("ExecExWhere returned nil")
	}

	// 删除操作
	result, err := updatedModel.Delete()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
