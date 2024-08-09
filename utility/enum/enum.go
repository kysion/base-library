package enum

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type IEnumCodeInt interface {
	Code() int
	ToMap() map[string]any
	// Description returns the brief description for current code.
	Description() string
}
type IEnumCodeStr interface {
	Code() string
	ToMap() map[string]any
	// Description returns the brief description for current code.
	Description() string
}

type IEnumCode[TCode uint | uint8 | uint16 | uint32 | uintptr | uint64 | int | int8 | int16 | int32 | int64 | string] interface {
	Code() TCode
	ToMap() map[string]any
	// Description returns the brief description for current code.
	Description() string
	// Has 是否有 enumType, 多个则全部包含返回 true
	Has(enumType ...IEnumCode[TCode]) bool
	// Add 自减
	Add(enumType ...IEnumCode[TCode]) bool
	// Remove 自加
	Remove(enumType ...IEnumCode[TCode]) bool
}

type IEnumCodeWithData[TCode uint | uint8 | uint16 | uint32 | uintptr | uint64 | int | int8 | int16 | int32 | int64 | string, TData any] interface {
	Code() TCode
	Data() TData
	ToMap() map[string]any
	Description() string
	// Has 是否有 enumType, 多个则全部包含返回 true
	Has(enumInfo ...IEnumCode[TCode]) bool
	// Add 自减
	Add(enumInfo ...IEnumCode[TCode]) bool
	// Remove 自加
	Remove(enumInfo ...IEnumCode[TCode]) bool
}

// EnumType [T any] is an implementer for interface Code for internal usage only.
type enumType[TCode uint | uint8 | uint16 | uint32 | uintptr | uint64 | int | int8 | int16 | int32 | int64 | string, TData any] struct {
	code        TCode  // Error code, usually an integer.
	data        TData  // Brief data for this value.
	description string // Brief description for this code.
}

// Code returns the integer number of current code.
func (e *enumType[TCode, TData]) Code() TCode {
	return e.code
}

// Description returns the brief description for current code.
func (e *enumType[TCode, TData]) Description() string {
	return e.description
}

// Data returns the T data of current code.
func (e *enumType[TCode, TData]) Data() TData {
	return e.data
}

func (e *enumType[TCode, TData]) ToMap() map[string]any {
	return map[string]any{
		"code":        e.code,
		"description": e.description,
		"data":        e.data,
	}
}

// Has 是否有 enumType, 多个则全部包含返回 true
func (e *enumType[TCode, TData]) Has(enumInfo ...IEnumCode[TCode]) bool {
	if len(enumInfo) <= 0 {
		return false
	}

	for _, item := range enumInfo {
		var codeAny any = item.Code()

		_, ok := codeAny.(string)

		if ok && item.Code() != e.code {
			return false
		}

		localCode := gconv.Int64(item.Code())
		if gconv.Int64(e.code)&localCode != localCode {
			return false
		}
	}

	return true
}

// Add 自加
func (e *enumType[TCode, TData]) Add(enumInfo ...IEnumCode[TCode]) bool {
	if len(enumInfo) <= 0 {
		return false
	}

	var codeAny any = e.code

	_, ok := codeAny.(string)

	if ok {
		return false
	}

	newCode := gconv.Int64(e.code)

	for _, item := range enumInfo {
		localCode := gconv.Int64(item.Code())
		if gconv.Int64(e.code)&localCode != localCode {
			newCode = newCode | localCode
		}
	}

	if e.code == TCode(newCode) {
		return false
	}

	e.code = TCode(newCode)

	return true
}

// Remove 自减
func (e *enumType[TCode, TData]) Remove(enumInfo ...IEnumCode[TCode]) bool {
	if len(enumInfo) <= 0 {
		return false
	}

	var codeAny any = e.code

	_, ok := codeAny.(string)

	if ok {
		return false
	}

	newCode := gconv.Int64(e.code)

	for _, item := range enumInfo {
		localCode := gconv.Int64(item.Code())
		if gconv.Int64(e.code)&localCode != localCode {
			newCode = newCode & ^localCode
		}
	}

	if e.code == TCode(newCode) {
		return false
	}

	e.code = TCode(newCode)
	return true
}

func New[R IEnumCode[TCode], TCode uint | uint8 | uint16 | uint32 | uintptr | uint64 | int | int8 | int16 | int32 | int64 | string](code TCode, description string) R {
	var result interface{}
	result = &enumType[TCode, interface{}]{
		code:        code,
		description: description,
	}
	return result.(R)
}

func NewWithData[TCode uint | uint8 | uint16 | uint32 | uintptr | uint64 | int | int8 | int16 | int32 | int64 | string, TData any](code TCode, data TData, description string) *IEnumCodeWithData[TCode, TData] {
	var result interface{}
	result = &enumType[TCode, TData]{
		code:        code,
		data:        data,
		description: description,
	}
	return result.(*IEnumCodeWithData[TCode, TData])
}

func GetTypes[V uint | uint8 | uint16 | uint32 | uintptr | uint64 | int | int8 | int16 | int32 | int64, T IEnumCode[V]](code V, enumOjb interface{}) []IEnumCode[V] {
	typeMaps := gconv.Map(enumOjb)

	result := make([]IEnumCode[V], 0)

	for key, value := range typeMaps {
		var item = typeMaps[key].(IEnumCode[V])
		if code&item.Code() == item.Code() {
			result = append(result, item)
		}

		g.Dump(key, value, item.Code())
	}

	return result
}
