package enum_benchmark

import (
	"testing"
)

// SimpleIntEnum 是一个简单的整型枚举实现
type SimpleIntEnum struct {
	code        int
	description string
}

func (e *SimpleIntEnum) Code() int {
	return e.code
}

func (e *SimpleIntEnum) Description() string {
	return e.description
}

func (e *SimpleIntEnum) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"code":        e.code,
		"description": e.description,
	}
}

// Has 检查是否有指定的枚举类型
func (e *SimpleIntEnum) Has(enumInfo ...*SimpleIntEnum) bool {
	if len(enumInfo) <= 0 {
		return false
	}

	for _, item := range enumInfo {
		if e.code&item.code != item.code {
			return false
		}
	}

	return true
}

// Add 添加指定的枚举类型
func (e *SimpleIntEnum) Add(enumInfo ...*SimpleIntEnum) bool {
	if len(enumInfo) <= 0 {
		return false
	}

	newCode := e.code

	for _, item := range enumInfo {
		newCode |= item.code
	}

	if e.code == newCode {
		return false
	}

	e.code = newCode
	return true
}

// Remove 移除指定的枚举类型
func (e *SimpleIntEnum) Remove(enumInfo ...*SimpleIntEnum) bool {
	if len(enumInfo) <= 0 {
		return false
	}

	newCode := e.code

	for _, item := range enumInfo {
		newCode &= ^item.code
	}

	if e.code == newCode {
		return false
	}

	e.code = newCode
	return true
}

// 预定义一些测试用的枚举值
var (
	SEnum1 = &SimpleIntEnum{code: 1, description: "测试枚举1"}
	SEnum2 = &SimpleIntEnum{code: 2, description: "测试枚举2"}
	SEnum4 = &SimpleIntEnum{code: 4, description: "测试枚举4"}
	SEnum8 = &SimpleIntEnum{code: 8, description: "测试枚举8"}
)

// BenchmarkNewSimpleEnum 测试创建简单枚举的性能
func BenchmarkNewSimpleEnum(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = &SimpleIntEnum{code: i, description: "测试枚举"}
	}
}

// BenchmarkSimpleCode 测试获取枚举代码的性能
func BenchmarkSimpleCode(b *testing.B) {
	enum := SEnum1
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = enum.Code()
	}
}

// BenchmarkSimpleDescription 测试获取枚举描述的性能
func BenchmarkSimpleDescription(b *testing.B) {
	enum := SEnum1
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = enum.Description()
	}
}

// BenchmarkSimpleToMap 测试枚举转换为Map的性能
func BenchmarkSimpleToMap(b *testing.B) {
	enum := SEnum1
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = enum.ToMap()
	}
}

// BenchmarkSimpleHasSingleEnum 测试枚举Has方法(单个枚举)的性能
func BenchmarkSimpleHasSingleEnum(b *testing.B) {
	enum := SEnum1
	checkEnum := SEnum1
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = enum.Has(checkEnum)
	}
}

// BenchmarkSimpleHasMultipleEnums 测试枚举Has方法(多个枚举)的性能
func BenchmarkSimpleHasMultipleEnums(b *testing.B) {
	enum := &SimpleIntEnum{code: 15, description: "组合枚举"} // 1 + 2 + 4 + 8 = 15
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = enum.Has(SEnum1, SEnum2, SEnum4, SEnum8)
	}
}

// BenchmarkSimpleAddSingleEnum 测试枚举Add方法(单个枚举)的性能
func BenchmarkSimpleAddSingleEnum(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		enum := &SimpleIntEnum{code: 1, description: "测试枚举"}
		_ = enum.Add(SEnum2)
	}
}

// BenchmarkSimpleAddMultipleEnums 测试枚举Add方法(多个枚举)的性能
func BenchmarkSimpleAddMultipleEnums(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		enum := &SimpleIntEnum{code: 1, description: "测试枚举"}
		_ = enum.Add(SEnum2, SEnum4, SEnum8)
	}
}

// BenchmarkSimpleRemoveSingleEnum 测试枚举Remove方法(单个枚举)的性能
func BenchmarkSimpleRemoveSingleEnum(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		enum := &SimpleIntEnum{code: 3, description: "测试枚举"} // 1 + 2 = 3
		_ = enum.Remove(SEnum1)
	}
}

// BenchmarkSimpleRemoveMultipleEnums 测试枚举Remove方法(多个枚举)的性能
func BenchmarkSimpleRemoveMultipleEnums(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		enum := &SimpleIntEnum{code: 15, description: "测试枚举"} // 1 + 2 + 4 + 8 = 15
		_ = enum.Remove(SEnum1, SEnum2)
	}
}

// BenchmarkSimpleParallelOperations 测试并发操作枚举的性能
func BenchmarkSimpleParallelOperations(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			enum := &SimpleIntEnum{code: 1, description: "测试枚举"}
			_ = enum.Add(SEnum2)
			_ = enum.Has(SEnum1)
			_ = enum.Has(SEnum2)
			_ = enum.Remove(SEnum1)
			_ = enum.ToMap()
		}
	})
}
