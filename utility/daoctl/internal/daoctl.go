package internal

import (
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/kysion/base-library/base_model"
)

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
	if len(orderBy) == 0 {
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
	if len(searchFieldArr) == 0 {
		return db, nil
	}

	builder := db.Builder()

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
					builder = builder.WhereOrNotIn(field.Field, field.Value)
				} else {
					builder = builder.WhereNotIn(field.Field, field.Value)
				}
			} else {
				if field.IsOrWhere {
					builder = builder.WhereOrIn(field.Field, field.Value)
				} else {
					builder = builder.WhereIn(field.Field, field.Value)
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
					builder = builder.WhereOrNotBetween(field.Field, minValue, maxValue)
				} else {
					builder = builder.WhereNotBetween(field.Field, minValue, maxValue)
				}
			} else {
				if field.IsOrWhere {
					builder = builder.WhereOrBetween(field.Field, minValue, maxValue)
				} else {
					builder = builder.WhereBetween(field.Field, minValue, maxValue)
				}
			}

		case whereClause == "like":
			// 处理LIKE查询条件
			if modifierClause == "not" {
				if field.IsOrWhere {
					builder = builder.WhereOrNotLike(field.Field, field.Value)
				} else {
					builder = builder.WhereNotLike(field.Field, field.Value)
				}
			} else {
				if field.IsOrWhere {
					builder = builder.WhereOrLike(field.Field, field.Value)
				} else {
					builder = builder.WhereLike(field.Field, gconv.String(field.Value))
				}
			}

		default:
			// 处理其他查询条件，如>、<、=等
			if gstr.Contains(field.Field, "&") {
				builder = builder.Wheref(field.Field+" "+field.Where+" ?", gconv.String(field.Value))
			} else {
				// 使用映射表简化代码逻辑
				switch field.Where {
				case ">":
					builder = applyComparisonOperator(builder, field, "GT")
				case ">=":
					builder = applyComparisonOperator(builder, field, "GTE")
				case "<":
					builder = applyComparisonOperator(builder, field, "LT")
				case "<=":
					builder = applyComparisonOperator(builder, field, "LTE")
				case "<>":
					if field.IsOrWhere {
						builder = builder.WhereOrNotIn(field.Field, field.Value)
					} else {
						builder = builder.WhereNotIn(field.Field, field.Value)
					}
				case "=":
					if field.IsOrWhere {
						builder = builder.WhereOr(field.Field, field.Value)
					} else {
						builder = builder.Where(field.Field, field.Value)
					}
				default:
					// 如果查询操作符不支持，则返回错误
					return nil, gerror.New("查询条件参数错误")
				}
			}
		}

		if len(field.Children) > 0 {
			if field.IsOrWhere {
				return MakeBuilder(db.WhereOr(builder), field.Children)
			} else {
				return MakeBuilder(db.Where(builder), field.Children)
			}
		}
	}
	// 返回构建完成的数据库模型对象
	return db.Where(builder), nil
}

// applyComparisonOperator 辅助函数，用于应用比较操作符
// 参数:
// - db: 数据库模型
// - field: 过滤字段信息
// - operator: 操作符类型 (GT, GTE, LT, LTE)
// 返回值:
// - 更新后的数据库模型
func applyComparisonOperator(db *gdb.WhereBuilder, field base_model.FilterInfo, operator string) *gdb.WhereBuilder {
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
