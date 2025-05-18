package daoctl

import (
	"math"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/kysion/base-library/base_model"
	"github.com/kysion/base-library/utility/daoctl/internal"
)

// GetById 根据ID获取模型对象。
//
// 该函数通过查询指定ID的记录，返回相应的模型实例。
// 使用了泛型T，使得该函数可以返回任何类型的模型实例。
//
// 参数:
//
//	model: *gdb.Model 类型，表示数据库模型，用于执行查询。
//	id: int64 类型，表示要查询的模型的唯一标识符。
//
// 返回值:
//
//	*T: 返回一个与模型相关联的实例，T为泛型，可以是任何类型。
//
// 为什么这么做:
//
//	该函数提供了一种通用的方法来根据ID查询数据库，并将结果转换为指定的模型类型，
//	从而避免了重复的类型断言和降低了代码的复杂度。
func GetById[T any](model *gdb.Model, id int64) *T {
	return Scan[T](model.Where("id", id))
}

// GetByIdWithError 根据给定的ID获取模型对象。
// 该函数使用了泛型T，允许返回任何类型的模型实例。
// 参数model指定了要查询的模型，参数id为要查询的模型的唯一标识符。
// 函数返回一个*T类型的指针和一个可能的错误。
// 如果查询成功，将返回模型的实例和nil错误；如果查询失败，将返回nil和错误详情。
func GetByIdWithError[T any](model *gdb.Model, id int64) (*T, error) {
	return ScanWithError[T](model.Where("id", id))
}

// Find 根据指定条件查找数据，并按指定字段排序。
//
// 参数:
// - model: 数据模型，用于指定查询的表。
// - orderBy: 排序条件数组，用于指定查询结果的排序顺序。
// - searchFields: 查询条件数组，用于指定查询时要匹配的字段和值。
//
// 返回值:
// - response: 包含查询结果的结构体指针，其中包含符合查询条件的数据集合。
// - err: 错误信息，如果查询过程中发生错误，则返回该错误。
func Find[T any](model *gdb.Model, orderBy []base_model.OrderBy, searchFields ...base_model.FilterInfo) (response *base_model.CollectRes[T], err error) {
	// ExecExWhere 执行动态的 where 条件查询，用于在查询前对 model 进行筛选。
	model = ExecExWhere(model)

	// 使用 Query 函数执行实际的查询操作，传入的参数定义了查询的详细要求，如过滤条件、分页和排序。
	// 这里将查询设置为从第一页开始，且不进行分页（PageSize为-1）。
	return Query[T](model, &base_model.SearchParams{
		Filter: searchFields,
		Pagination: base_model.Pagination{
			PageNum:  1,
			PageSize: -1,
		},
		OrderBy: orderBy,
	}, true)
}

// GetAll 从数据库中获取所有符合条件的实体。
//
// 参数：
// - model: GORM模型，用于指定查询的模型。
// - info: 分页信息结构，包含分页参数。如果为nil，将返回所有数据。
//
// 返回值：
// - response: 包含查询结果和分页信息的结构。
// - err: 查询过程中可能发生的错误。
func GetAll[T any](model *gdb.Model, info *base_model.Pagination) (response *base_model.CollectRes[*T], err error) {
	// 对模型应用额外的查询条件。
	model = ExecExWhere(model)

	// 计算满足条件的总记录数。
	total, err := model.Count()
	// 初始化实体切片，预设容量为总数。
	entities := make([]*T, 0, total)
	// 如果没有提供分页信息，默认返回所有数据。
	if info == nil {
		info = &base_model.Pagination{
			PageNum:  1,
			PageSize: gconv.Int(total),
		}
	}

	// 检查计数操作是否有错误。
	if err != nil {
		return
	}
	// 执行分页查询。
	err = model.Page(info.PageNum, info.PageSize).Scan(&entities)

	// 构造并返回结果集和分页信息。
	return &base_model.CollectRes[*T]{
		Records: entities,
		PaginationRes: base_model.PaginationRes{
			Pagination: *info,
			PageTotal:  gconv.Int(math.Ceil(gconv.Float64(total) / gconv.Float64(info.PageSize))),
			Total:      gconv.Int64(total),
		},
	}, nil
}

// Query 函数用于执行数据库查询操作，并返回查询结果。
// [T] 是一个泛型参数，允许函数处理各种类型的查询结果。
// 参数:
// - model: 指向数据库模型的指针，用于指定查询的数据库表。
// - searchFields: 指向搜索参数的指针，包含过滤和排序等信息。如果为nil，将使用默认搜索参数。
// - IsExport: 一个布尔值，指示是否为导出操作。如果是导出操作，将返回所有记录，而不是分页数据。
// 返回值:
// - response: 包含查询结果和分页信息的指针。
// - err: 执行查询过程中可能发生的错误。
func Query[T any](model *gdb.Model, searchFields *base_model.SearchParams, IsExport bool) (*base_model.CollectRes[T], error) {
	// 对模型执行预处理，可能包括设置默认的查询条件等。
	model = ExecExWhere(model)

	// 如果没有提供搜索参数，则初始化一个默认的搜索参数对象。
	if searchFields == nil {
		searchFields = &base_model.SearchParams{}
	}

	// 根据过滤条件构建查询语句。
	queryDb, _ := internal.MakeBuilder(model, searchFields.Filter)
	// 根据排序条件应用排序。
	queryDb = internal.MakeOrderBy(queryDb, searchFields.OrderBy)

	// 确保页码至少为1，防止无效的页码值。
	if searchFields.PageNum <= 1 {
		searchFields.PageNum = 1
	}

	// 确保页大小至少为正数，如果为0或负数，则设置为默认值20。
	if searchFields.PageSize <= 0 {
		searchFields.PageSize = 20
	}

	var err error
	count := 0

	// 初始化一个空的实体切片，用于存储查询结果。
	entities := make([]T, 0)
	// 如果是导出操作，设置页大小为-1，以获取所有记录。
	if IsExport {
		// 执行查询并存储结果到entities切片中。
		err = queryDb.ScanAndCount(&entities, &count, false)
	} else {
		// 执行分页查询，并存储结果到entities切片中。
		err = queryDb.Page(searchFields.PageNum, searchFields.PageSize).ScanAndCount(&entities, &count, false)
	}

	response := &base_model.CollectRes[T]{
		Records: entities,
		PaginationRes: base_model.PaginationRes{
			Pagination: base_model.Pagination{
				PageNum:  searchFields.PageNum,
				PageSize: searchFields.PageSize,
			},
			Total:     int64(count),
			PageTotal: gconv.Int(math.Ceil(gconv.Float64(count) / gconv.Float64(searchFields.PageSize))),
		},
	}

	return response, err
}

// MakeModel 函数用于创建一个查询模型，并返回一个指向该模型的指针。
func MakeModel(model *gdb.Model, searchFields *base_model.SearchParams) *gdb.Model {
	// 对模型执行预处理，可能包括设置默认的查询条件等。
	model = ExecExWhere(model)

	// 如果没有提供搜索参数，则初始化一个默认的搜索参数对象。
	if searchFields == nil {
		searchFields = &base_model.SearchParams{}
	}

	// 根据过滤条件构建查询语句。
	queryDb, _ := internal.MakeBuilder(model, searchFields.Filter)
	// 根据排序条件应用排序。
	queryDb = internal.MakeOrderBy(queryDb, searchFields.OrderBy)

	return queryDb
}
