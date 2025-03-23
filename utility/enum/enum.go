package enum

import (
	"github.com/gogf/gf/v2/util/gconv"
)

// IEnumCodeInt 是一个接口，用于处理整型枚举代码。
type IEnumCodeInt interface {
	Code() int
	ToMap() map[string]any
	// Description 返回当前代码的简短描述。
	Description() string
}

// IEnumCodeStr 是一个接口，用于处理字符串型枚举代码。
type IEnumCodeStr interface {
	Code() string
	ToMap() map[string]any
	// Description 返回当前代码的简短描述。
	Description() string
}

// IEnumCode 是一个泛型接口，用于处理不同类型的枚举代码。
type IEnumCode[TCode uint | uint8 | uint16 | uint32 | uintptr | uint64 | int | int8 | int16 | int32 | int64 | string] interface {
	Code() TCode
	ToMap() map[string]any
	// Description 返回当前代码的简短描述。
	Description() string
	// Has 检查是否有指定的枚举类型，如果有多个，则必须全部包含才返回 true。
	Has(enumType ...IEnumCode[TCode]) bool
	// Add 添加指定的枚举类型。
	Add(enumType ...IEnumCode[TCode]) bool
	// Remove 移除指定的枚举类型。
	Remove(enumType ...IEnumCode[TCode]) bool
}

// IEnumCodeWithData 是一个泛型接口，用于处理带有附加数据的枚举代码。
type IEnumCodeWithData[TCode uint | uint8 | uint16 | uint32 | uintptr | uint64 | int | int8 | int16 | int32 | int64 | string, TData any] interface {
	Code() TCode
	Data() TData
	ToMap() map[string]any
	Description() string
	// Has 检查是否有指定的枚举类型，如果有多个，则必须全部包含才返回 true。
	Has(enumInfo ...IEnumCode[TCode]) bool
	// Add 添加指定的枚举类型。
	Add(enumInfo ...IEnumCode[TCode]) bool
	// Remove 移除指定的枚举类型。
	Remove(enumInfo ...IEnumCode[TCode]) bool
}

// enumType 是一个内部使用的实现，为接口 Code 提供具体实现。
type enumType[TCode uint | uint8 | uint16 | uint32 | uintptr | uint64 | int | int8 | int16 | int32 | int64 | string, TData any] struct {
	code        TCode  // 错误代码，通常是一个整数。
	data        TData  // 该值的简短数据。
	description string // 该代码的简短描述。
}

// Code 返回当前代码的数值。
func (e *enumType[TCode, TData]) Code() TCode {
	return e.code
}

// Description 返回当前代码的简短描述。
func (e *enumType[TCode, TData]) Description() string {
	return e.description
}

// Data 返回当前代码的附加数据。
func (e *enumType[TCode, TData]) Data() TData {
	return e.data
}

// ToMap 将枚举类型转换为映射格式。
func (e *enumType[TCode, TData]) ToMap() map[string]any {
	return map[string]any{
		"code":        e.code,
		"description": e.description,
		"data":        e.data,
	}
}

// Has 检查是否有指定的枚举类型，如果有多个，则必须全部包含才返回 true。
func (e *enumType[TCode, TData]) Has(enumInfo ...IEnumCode[TCode]) bool {
	if len(enumInfo) <= 0 {
		return false
	}

	for _, item := range enumInfo {
		if _, ok := any(e.code).(string); ok { // 检查是否为字符串类型
			if e.code != item.Code() {
				return false
			}
		} else { // 数值类型
			switch code := any(e.code).(type) {
			case int:
				itemCode := any(item.Code()).(int)
				if code&itemCode != itemCode {
					return false
				}
			case int8:
				itemCode := any(item.Code()).(int8)
				if code&itemCode != itemCode {
					return false
				}
			case int16:
				itemCode := any(item.Code()).(int16)
				if code&itemCode != itemCode {
					return false
				}
			case int32:
				itemCode := any(item.Code()).(int32)
				if code&itemCode != itemCode {
					return false
				}
			case int64:
				itemCode := any(item.Code()).(int64)
				if code&itemCode != itemCode {
					return false
				}
			case uint:
				itemCode := any(item.Code()).(uint)
				if code&itemCode != itemCode {
					return false
				}
			case uint8:
				itemCode := any(item.Code()).(uint8)
				if code&itemCode != itemCode {
					return false
				}
			case uint16:
				itemCode := any(item.Code()).(uint16)
				if code&itemCode != itemCode {
					return false
				}
			case uint32:
				itemCode := any(item.Code()).(uint32)
				if code&itemCode != itemCode {
					return false
				}
			case uint64:
				itemCode := any(item.Code()).(uint64)
				if code&itemCode != itemCode {
					return false
				}
			default:
				return false
			}
		}
	}

	return true
}

// Add 添加指定的枚举类型。
func (e *enumType[TCode, TData]) Add(enumInfo ...IEnumCode[TCode]) bool {
	if len(enumInfo) <= 0 {
		return false
	}

	// 字符串类型不支持位运算
	if _, ok := any(e.code).(string); ok {
		return false
	}

	var changed bool
	var oldCode = e.code

	// 从Go 1.18开始支持泛型，但泛型类型断言还不完善
	// 这里使用更简单的方式实现
	for _, item := range enumInfo {
		switch code := any(e.code).(type) {
		case int:
			itemCode := any(item.Code()).(int)
			if code|itemCode != code {
				e.code = TCode(code | itemCode)
				changed = true
			}
		case int8:
			itemCode := any(item.Code()).(int8)
			if code|itemCode != code {
				e.code = TCode(code | itemCode)
				changed = true
			}
		case int16:
			itemCode := any(item.Code()).(int16)
			if code|itemCode != code {
				e.code = TCode(code | itemCode)
				changed = true
			}
		case int32:
			itemCode := any(item.Code()).(int32)
			if code|itemCode != code {
				e.code = TCode(code | itemCode)
				changed = true
			}
		case int64:
			itemCode := any(item.Code()).(int64)
			if code|itemCode != code {
				e.code = TCode(code | itemCode)
				changed = true
			}
		case uint:
			itemCode := any(item.Code()).(uint)
			if code|itemCode != code {
				e.code = TCode(code | itemCode)
				changed = true
			}
		case uint8:
			itemCode := any(item.Code()).(uint8)
			if code|itemCode != code {
				e.code = TCode(code | itemCode)
				changed = true
			}
		case uint16:
			itemCode := any(item.Code()).(uint16)
			if code|itemCode != code {
				e.code = TCode(code | itemCode)
				changed = true
			}
		case uint32:
			itemCode := any(item.Code()).(uint32)
			if code|itemCode != code {
				e.code = TCode(code | itemCode)
				changed = true
			}
		case uint64:
			itemCode := any(item.Code()).(uint64)
			if code|itemCode != code {
				e.code = TCode(code | itemCode)
				changed = true
			}
		}
	}

	// 对于性能测试，如果出错则恢复原值
	if !changed {
		e.code = oldCode
	}

	return changed
}

// Remove 移除指定的枚举类型。
func (e *enumType[TCode, TData]) Remove(enumInfo ...IEnumCode[TCode]) bool {
	if len(enumInfo) <= 0 {
		return false
	}

	// 字符串类型不支持位运算
	if _, ok := any(e.code).(string); ok {
		return false
	}

	var changed bool
	var oldCode = e.code

	// 从Go 1.18开始支持泛型，但泛型类型断言还不完善
	// 这里使用更简单的方式实现
	for _, item := range enumInfo {
		switch code := any(e.code).(type) {
		case int:
			itemCode := any(item.Code()).(int)
			if code&itemCode != 0 {
				e.code = TCode(code &^ itemCode)
				changed = true
			}
		case int8:
			itemCode := any(item.Code()).(int8)
			if code&itemCode != 0 {
				e.code = TCode(code &^ itemCode)
				changed = true
			}
		case int16:
			itemCode := any(item.Code()).(int16)
			if code&itemCode != 0 {
				e.code = TCode(code &^ itemCode)
				changed = true
			}
		case int32:
			itemCode := any(item.Code()).(int32)
			if code&itemCode != 0 {
				e.code = TCode(code &^ itemCode)
				changed = true
			}
		case int64:
			itemCode := any(item.Code()).(int64)
			if code&itemCode != 0 {
				e.code = TCode(code &^ itemCode)
				changed = true
			}
		case uint:
			itemCode := any(item.Code()).(uint)
			if code&itemCode != 0 {
				e.code = TCode(code &^ itemCode)
				changed = true
			}
		case uint8:
			itemCode := any(item.Code()).(uint8)
			if code&itemCode != 0 {
				e.code = TCode(code &^ itemCode)
				changed = true
			}
		case uint16:
			itemCode := any(item.Code()).(uint16)
			if code&itemCode != 0 {
				e.code = TCode(code &^ itemCode)
				changed = true
			}
		case uint32:
			itemCode := any(item.Code()).(uint32)
			if code&itemCode != 0 {
				e.code = TCode(code &^ itemCode)
				changed = true
			}
		case uint64:
			itemCode := any(item.Code()).(uint64)
			if code&itemCode != 0 {
				e.code = TCode(code &^ itemCode)
				changed = true
			}
		}
	}

	// 对于性能测试，如果出错则恢复原值
	if !changed {
		e.code = oldCode
	}

	return changed
}

// New 创建一个新的枚举类型实例。
func New[R IEnumCode[TCode], TCode uint | uint8 | uint16 | uint32 | uintptr | uint64 | int | int8 | int16 | int32 | int64 | string](code TCode, description string) R {
	var result interface{}
	result = &enumType[TCode, interface{}]{
		code:        code,
		description: description,
	}
	return result.(R)
}

// NewWithData 创建一个新的带有附加数据的枚举类型实例。
func NewWithData[TCode uint | uint8 | uint16 | uint32 | uintptr | uint64 | int | int8 | int16 | int32 | int64 | string, TData any](code TCode, data TData, description string) *IEnumCodeWithData[TCode, TData] {
	var result interface{}
	result = &enumType[TCode, TData]{
		code:        code,
		data:        data,
		description: description,
	}
	return result.(*IEnumCodeWithData[TCode, TData])
}

// GetTypes 获取指定代码类型的所有枚举类型实例。
func GetTypes[V uint | uint8 | uint16 | uint32 | uintptr | uint64 | int | int8 | int16 | int32 | int64, T IEnumCode[V]](code V, enumOjb interface{}) []IEnumCode[V] {
	typeMaps := gconv.Map(enumOjb)

	result := make([]IEnumCode[V], 0)

	for key, _ := range typeMaps {
		var item = typeMaps[key].(IEnumCode[V])
		if code&item.Code() == item.Code() {
			result = append(result, item)
		}
	}

	return result
}
