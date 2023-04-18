package base_funs

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/kysion/base-library/base_model"
	"reflect"
)

func If[R any](condition bool, trueVal, falseVal R) R {
	if condition {
		return trueVal
	} else {
		return falseVal
	}
}

func SearchFilterEx(search base_model.SearchParams, fields ...string) *base_model.SearchParams {
	result := &base_model.SearchParams{}
	newFilter := make([]base_model.FilterInfo, 0)
	for _, info := range search.Filter {
		count := len(result.Filter)
		for _, field := range fields {
			if gstr.ToLower(info.Field) == gstr.ToLower(field) {
				newFilter = append(result.Filter, info)
			}
		}
		if count == len(result.Filter) {
			newFilter = append(newFilter, info)
		}
	}
	search.Filter = newFilter
	return result
}

// ByteCountSI 以1000作为基数
func ByteCountSI[T int64 | uint64](b T) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

// ByteCountIEC 以1024作为基数
func ByteCountIEC[T int64 | uint64](b T) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

func RemoveSliceAt[T int | int64 | string | uint | uint64](slice []T, elem T) []T {
	if len(slice) == 0 {
		return slice
	}

	for i, v := range slice {
		if v == elem {
			slice = append(slice[:i], slice[i+1:]...)
			return RemoveSliceAt(slice, elem)
			break
		}
	}
	return slice
}

func AttrBuilder[T any, TP any](ctx context.Context, key string, builder ...func(data TP)) context.Context {
	key = key + "::" + reflect.ValueOf(new(T)).Type().String() + "::" + reflect.ValueOf(new(TP)).Type().String()
	key = gstr.Replace(key, "*", "")
	def := func(data TP) {}

	if len(builder) > 0 {
		def = builder[0]
	}

	return context.WithValue(ctx, key,
		base_model.KeyValueT[string, func(data TP)]{
			Key:   key,
			Value: def,
		},
	)
}

//union_main_id::co_model.EmployeeRes::[]co_model.Team
func AttrMake[T any, TP any](ctx context.Context, key string, builder func() TP) {
	key = key + "::" + reflect.ValueOf(new(T)).Type().String() + "::" + reflect.ValueOf(new(TP)).Type().String()
	key = gstr.Replace(key, "*", "")
	v := ctx.Value(key)

	data, has := v.(base_model.KeyValueT[string, func(data TP)])
	if v != nil && has {
		data.Value(builder())
	}
}
