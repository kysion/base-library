package daoctl

import (
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/kysion/base-library/base_model"
	"github.com/kysion/base-library/utility/daoctl/internal"
	"math"
)

func GetById[T any](model *gdb.Model, id int64) *T {
	return Scan[T](model.Where("id", id))
}

func GetByIdWithError[T any](model *gdb.Model, id int64) (*T, error) {
	return ScanWithError[T](model.Where("id", id))
}

func Find[T any](db *gdb.Model, orderBy []base_model.OrderBy, searchFields ...base_model.FilterInfo) (response *base_model.CollectRes[T], err error) {
	return Query[T](db, &base_model.SearchParams{
		Filter: searchFields,
		Pagination: base_model.Pagination{
			PageNum:  1,
			PageSize: -1,
		},
		OrderBy: orderBy,
	}, true)
}

func GetAll[T any](db *gdb.Model, info *base_model.Pagination) (response *base_model.CollectRes[*T], err error) {
	total, err := db.Count()
	entities := make([]*T, 0, total)
	if info == nil {
		info = &base_model.Pagination{
			PageNum:  1,
			PageSize: gconv.Int(total),
		}
	}

	if err != nil {
		return
	}
	err = db.Page(info.PageNum, info.PageSize).Scan(&entities)

	return &base_model.CollectRes[*T]{
		Records: entities,
		PaginationRes: base_model.PaginationRes{
			Pagination: *info,
			PageTotal:  gconv.Int(math.Ceil(gconv.Float64(total) / gconv.Float64(info.PageSize))),
			Total:      gconv.Int64(total),
		},
	}, nil
}

func Query[T any](db *gdb.Model, searchFields *base_model.SearchParams, IsExport bool) (response *base_model.CollectRes[T], err error) {
	// 查询具体的值
	queryDb, _ := internal.MakeBuilder(db, searchFields.Filter)
	queryDb = internal.MakeOrderBy(queryDb, searchFields.OrderBy)

	if searchFields.PageSize == 0 {
		searchFields.PageSize = 20
		searchFields.PageNum = 1
	}

	entities := make([]T, 0)
	if searchFields == nil || IsExport {
		searchFields.PageSize = -1
		err = queryDb.Scan(&entities)
	} else {
		err = queryDb.Page(searchFields.PageNum, searchFields.PageSize).Scan(&entities)
	}

	response = &base_model.CollectRes[T]{
		Records:       entities,
		PaginationRes: internal.MakePaginationArr(db, searchFields.Pagination, searchFields.Filter),
	}

	return response, nil
}
