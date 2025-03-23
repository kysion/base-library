package internal

import (
	"math"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/kysion/base-library/base_model"
)

// MakeCountArr 生成计数数组
// 该函数的目的是为了统计满足过滤条件的记录总数
// 参数:
//
//	db: 数据库模型，用于执行数据库操作
//	searchFields: 过滤条件数组，包含需要进行过滤的字段信息
//
// 返回值:
//
//	total: 满足条件的记录总数
func MakeCountArr(db *gdb.Model, searchFields []base_model.FilterInfo) (total int64) {
	// 使用MakeBuilder函数根据提供的数据库模型和过滤条件构建查询语句
	db, err := MakeBuilder(db, searchFields)
	if err != nil {
		// 如果构建查询语句过程中发生错误，直接返回0
		return 0
	}

	// 初始化计数变量
	count := 0
	// 执行查询并统计满足条件的记录数，不需实际加载数据
	err = db.ScanAndCount(nil, &count, false)
	if err != nil {
		// 如果查询执行出错，直接返回0
		return 0
	}

	// 将统计得到的总数转换为int64类型并返回
	return gconv.Int64(count)
}

// MakeOrderBy1 根据指定的排序条件更新数据库查询模型
// 该函数接收一个数据库模型和一个排序条件数组，根据数组中的排序条件
// 更新数据库查询模型的排序设置，以便在执行查询时按照指定的顺序排序结果。
// 参数:
// - db: *gdb.Model - 一个指向数据库模型的指针，该模型代表了数据库的查询状态。
// - orderBy: []base_model.OrderBy - 一个包含排序条件的切片，每个排序条件指定一个字段和一个排序方向（升序或降序）。
// 返回值:
// - *gdb.Model - 返回更新后的数据库查询模型指针。
// @deprecated 请使用 MakeOrderBy 函数替代
func MakeOrderBy1(db *gdb.Model, orderBy []base_model.OrderBy) *gdb.Model {
	// 为保持向后兼容性，调用 MakeOrderBy 函数
	return MakeOrderBy(db, orderBy)
}

// MakeOrderBy 根据指定的排序条件更新数据库查询模型
// 参数:
//
//	db: *gdb.Model - 初始的数据库查询模型
//	orderBy: []base_model.OrderBy - 包含排序字段和顺序的切片
//
// 返回值:
//
//	*gdb.Model - 更新了排序条件的数据库查询模型
func MakeOrderBy(db *gdb.Model, orderBy []base_model.OrderBy) *gdb.Model {
	// 检查 orderBy 是否为空且不为 nil
	if orderBy == nil || len(orderBy) == 0 {
		return db
	}

	// 遍历 orderBy 切片中的每个排序条件
	for _, orderField := range orderBy {
		// 将字段名转换为数据库查询所需的格式
		orderField.Field = gstr.CaseSnakeFirstUpper(orderField.Field)
		// 移除字段名中可能存在的引号，增加安全性
		orderField.Field = gstr.ReplaceIByMap(orderField.Field, map[string]string{"\"": "", "'": ""})

		// 将排序方向转换为小写进行比较
		sortLower := gstr.ToLower(orderField.Sort)
		// 根据排序方式更新数据库查询模型
		if sortLower == "asc" {
			// 如果排序方式为升序，则调用 OrderAsc 方法
			db = db.OrderAsc(orderField.Field)
		} else if sortLower == "desc" {
			// 如果排序方式为降序，则调用 OrderDesc 方法
			db = db.OrderDesc(orderField.Field)
		}
	}

	// 返回更新后的数据库查询模型
	return db
}

// MakeBuilder 根据提供的搜索字段数组构建数据库查询条件。
// 参数:
// - db: 初始的数据库模型对象。
// - searchFieldArr: 包含搜索字段信息的数组。
// 返回值:
// - 修改后的数据库模型对象。
// - 如果构建过程中出现错误，则返回错误信息。
func MakeBuilder(db *gdb.Model, searchFieldArr []base_model.FilterInfo) (*gdb.Model, error) {
	// 检查searchFieldArr是否为空，为空则直接返回原始模型
	if searchFieldArr == nil || len(searchFieldArr) == 0 {
		return db, nil
	}

	// 遍历searchFieldArr，对每个字段构建查询条件
	for index, field := range searchFieldArr {
		// 将字段名称转换为Snake Case格式，并首字母大写
		field.Field = gstr.CaseSnakeFirstUpper(field.Field)

		// 确保字段名称不为空
		if gconv.String(field.Field) == "" {
			continue
		}

		// 移除字段名称中的特殊字符，如引号，防止SQL注入
		field.Field = gstr.ReplaceIByMap(field.Field, map[string]string{"\"": "", "'": ""})

		// 第一个字段默认使用WHERE而不是OR WHERE
		if index == 0 {
			field.IsOrWhere = false
		}

		// 将查询条件转换为小写，便于比较
		whereClause := gstr.ToLower(field.Where)
		modifierClause := gstr.ToLower(field.Modifier)

		// 根据查询条件的类型执行相应的查询构建操作
		switch {
		case whereClause == "in":
			// 处理IN查询条件
			if modifierClause == "not" {
				if field.IsOrWhere {
					db = db.WhereOrNotIn(field.Field, field.Value)
				} else {
					db = db.WhereNotIn(field.Field, field.Value)
				}
			} else {
				if field.IsOrWhere {
					db = db.WhereOrIn(field.Field, field.Value)
				} else {
					db = db.WhereIn(field.Field, field.Value)
				}
			}

		case whereClause == "between":
			// 处理BETWEEN查询条件
			valueArr := gstr.SplitAndTrim(gconv.String(field.Value), ",")
			minValue := valueArr[0]
			maxValue := minValue
			if len(valueArr) > 1 {
				maxValue = valueArr[1]
			}

			if modifierClause == "not" {
				if field.IsOrWhere {
					db = db.WhereOrNotBetween(field.Field, minValue, maxValue)
				} else {
					db = db.WhereNotBetween(field.Field, minValue, maxValue)
				}
			} else {
				if field.IsOrWhere {
					db = db.WhereOrBetween(field.Field, minValue, maxValue)
				} else {
					db = db.WhereBetween(field.Field, minValue, maxValue)
				}
			}

		case whereClause == "like":
			// 处理LIKE查询条件
			if modifierClause == "not" {
				if field.IsOrWhere {
					db = db.WhereOrNotLike(field.Field, field.Value)
				} else {
					db = db.WhereNotLike(field.Field, field.Value)
				}
			} else {
				if field.IsOrWhere {
					db = db.WhereOrLike(field.Field, field.Value)
				} else {
					db = db.WhereLike(field.Field, gconv.String(field.Value))
				}
			}

		default:
			// 处理其他查询条件，如>、<、=等
			if gstr.Contains(field.Field, "&") {
				db = db.Wheref(field.Field+" "+field.Where+" ?", gconv.String(field.Value))
			} else {
				// 使用映射表简化代码逻辑
				switch field.Where {
				case ">":
					db = applyComparisonOperator(db, field, "GT")
				case ">=":
					db = applyComparisonOperator(db, field, "GTE")
				case "<":
					db = applyComparisonOperator(db, field, "LT")
				case "<=":
					db = applyComparisonOperator(db, field, "LTE")
				case "<>":
					if field.IsOrWhere {
						db = db.WhereOrNotIn(field.Field, field.Value)
					} else {
						db = db.WhereNotIn(field.Field, field.Value)
					}
				case "=":
					if field.IsOrWhere {
						db = db.WhereOr(field.Field, field.Value)
					} else {
						db = db.Where(field.Field, field.Value)
					}
				default:
					// 如果查询操作符不支持，则返回错误
					return nil, gerror.New("查询条件参数错误")
				}
			}
		}
	}

	// 返回构建完成的数据库模型对象
	return db, nil
}

// applyComparisonOperator 辅助函数，用于应用比较操作符
// 参数:
// - db: 数据库模型
// - field: 过滤字段信息
// - operator: 操作符类型 (GT, GTE, LT, LTE)
// 返回值:
// - 更新后的数据库模型
func applyComparisonOperator(db *gdb.Model, field base_model.FilterInfo, operator string) *gdb.Model {
	if field.IsOrWhere {
		switch operator {
		case "GT":
			return db.WhereOrGT(field.Field, field.Value)
		case "GTE":
			return db.WhereOrGTE(field.Field, field.Value)
		case "LT":
			return db.WhereOrLT(field.Field, field.Value)
		case "LTE":
			return db.WhereOrLTE(field.Field, field.Value)
		}
	} else {
		switch operator {
		case "GT":
			return db.WhereGT(field.Field, field.Value)
		case "GTE":
			return db.WhereGTE(field.Field, field.Value)
		case "LT":
			return db.WhereLT(field.Field, field.Value)
		case "LTE":
			return db.WhereLTE(field.Field, field.Value)
		}
	}
	return db
}

// MakePaginationArr 生成分页结果数组
// 参数:
// db: 数据库模型指针，用于执行查询
// pagination: 分页参数，包含页码和每页大小等信息
// searchFields: 用于搜索的字段数组，用于过滤查询结果
// 返回值:
// base_model.PaginationRes: 分页结果结构，包含分页信息和查询结果总数
func MakePaginationArr(db *gdb.Model, pagination base_model.Pagination, searchFields []base_model.FilterInfo) base_model.PaginationRes {
	// 计算满足条件的总记录数
	total := MakeCountArr(db, searchFields)

	// 准备返回的分页结果结构
	result := base_model.PaginationRes{
		Pagination: pagination,
		Total:      total,
	}

	// 如果总记录数为 0，直接返回空结果
	if total == 0 {
		result.PageTotal = 0
		return result
	}

	// 如果每页大小为 -1（表示不分页）
	if pagination.PageSize == -1 {
		// 将总记录数设置为每页大小
		result.PageSize = gconv.Int(total)

		// 如果总记录数为 0，设置默认每页大小为 20
		if result.PageSize == 0 {
			result.PageSize = 20
		}

		// 不分页时总页数为 1
		result.PageTotal = 1
		return result
	}

	// 计算总页数，向上取整
	result.PageTotal = gconv.Int(math.Ceil(gconv.Float64(total) / gconv.Float64(pagination.PageSize)))

	return result
}
