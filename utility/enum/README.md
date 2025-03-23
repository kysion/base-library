# Kysion 枚举工具包

此包提供了一个通用的枚举实现，支持整数和字符串类型的枚举，同时实现了位运算功能。基于Go 1.18+泛型实现，提供类型安全的枚举操作。

## 主要特性

- 支持整数和字符串枚举
- 提供位运算操作 (Has, Add, Remove)
- 支持枚举项附加数据
- 类型安全的实现，防止不安全的类型转换
- 高性能实现，基本操作零内存分配

## 类型安全的实现

本库使用Go的泛型和`any`类型断言机制，确保所有类型转换都是安全的，具体实现：

1. 使用泛型约束限制可用的枚举类型
2. 位运算通过`bitOr`和`bitAndNot`辅助函数执行
3. 在这些函数内使用`any`类型断言确保相同类型间的操作
4. 每个操作都验证类型并安全处理类型转换

这种实现方式避免了直接类型转换带来的潜在问题，同时保持了高性能。

## 使用示例

```go
// 整数枚举示例
type Permission int

const (
    PermissionRead Permission = 1 << iota
    PermissionWrite
    PermissionExecute
)

var e = enum.New(PermissionRead | PermissionWrite)

// 检查是否包含权限
if e.Has(PermissionRead) {
    // 有读权限
}

// 添加权限
e.Add(PermissionExecute)

// 移除权限
e.Remove(PermissionWrite)
```

```go
// 字符串枚举示例
type Role string

const (
    RoleAdmin Role = "admin"
    RoleUser  Role = "user"
    RoleGuest Role = "guest"
)

var r = enum.New(RoleAdmin)

// 检查角色
if r.Has(RoleAdmin) {
    // 是管理员
}

// 修改角色
r.Set(RoleUser)
```

## 性能测试

枚举工具包经过了严格的性能测试，以下是最新的性能测试结果摘要：

| 操作类型                | 操作耗时 (ns)   | 内存分配      |
|------------------------|----------------|--------------|
| 基本枚举访问操作        | < 0.2 ns       | 零内存分配    |
| 单枚举位运算操作        | < 0.8 ns       | 零内存分配    |
| 多枚举位运算操作        | < 1.5 ns       | 零内存分配    |
| ToMap转换操作          | ~ 25.7 ns      | 16B/操作     |
| 并行混合操作            | ~ 6.9 ns       | 16B/操作     |

完整的性能测试报告和分析请查看 [benchmark/README.md](benchmark/README.md)。

### 性能结论

- 枚举的基本操作（读取、位运算）性能极高，每秒可执行数十亿次操作
- 除ToMap操作外，大多数操作不产生内存分配，有利于减少GC压力
- 并发操作下表现稳定，适合高并发场景使用

## 版本历史

### v1.x.x (最新)

- 修复类型转换安全性问题
- 优化位运算操作实现
- 通过`any`类型断言确保类型安全

### v1.0.0

- 初始版本发布
- 支持整数和字符串枚举
- 实现位运算基本功能

## 目录结构

- `enum.go` - 枚举工具包的主要实现
- `benchmark/` - 性能测试相关文件
  - `enum_test.go` - 性能测试代码
  - `run.sh` - 运行性能测试的脚本
  - `README.md` - 性能测试文档和结果

## 实现说明

本工具包使用 Go 1.18+ 的泛型特性实现，提供了类型安全的枚举操作。主要包含：

- 通过 `NumberEnumCode` 接口约束支持的数值类型
- 通过辅助函数 `bitOr` 和 `bitAndNot` 安全处理位运算和类型转换
- 使用 `any` 和类型断言来确保类型安全性，避免不安全的类型转换

## 使用示例

### 整型枚举示例

整型枚举支持位运算，常用于权限、状态标志等场景：

```go
package main

import (
    "fmt"
    "github.com/kysion/base-library/utility/enum"
)

// 定义一个权限枚举接口
type Permission interface {
    enum.IEnumCode[int]
}

// 创建一些权限枚举值
var (
    ReadPerm = enum.New[Permission](1, "读取权限")
    WritePerm = enum.New[Permission](2, "写入权限")
    DeletePerm = enum.New[Permission](4, "删除权限")
    AdminPerm = enum.New[Permission](8, "管理权限")
)

func main() {
    // 创建一个包含读取和写入权限的组合
    userPerm := enum.New[Permission](0, "用户权限")
    userPerm.Add(ReadPerm, WritePerm)
    
    // 检查是否包含特定权限
    fmt.Println("有读取权限:", userPerm.Has(ReadPerm))     // true
    fmt.Println("有删除权限:", userPerm.Has(DeletePerm))   // false
    
    // 添加新权限
    userPerm.Add(DeletePerm)
    fmt.Println("添加后有删除权限:", userPerm.Has(DeletePerm)) // true
    
    // 移除权限
    userPerm.Remove(WritePerm)
    fmt.Println("移除后有写入权限:", userPerm.Has(WritePerm)) // false
}
```

### 字符串枚举示例

字符串枚举适用于枚举值需要更好可读性的场景，如API状态码、国家代码等：

```go
package main

import (
    "fmt"
    "github.com/kysion/base-library/utility/enum"
)

// 定义一个国家枚举接口
type Country interface {
    enum.IEnumCode[string]
}

// 创建一些国家枚举值
var (
    China = enum.New[Country]("CN", "中国")
    USA = enum.New[Country]("US", "美国")
    Japan = enum.New[Country]("JP", "日本")
    UK = enum.New[Country]("UK", "英国")
)

func main() {
    // 创建一个国家枚举
    country := enum.New[Country]("CN", "中国")
    
    // 检查是否为特定国家
    fmt.Println("是中国:", country.Has(China))     // true
    fmt.Println("是美国:", country.Has(USA))       // false
    
    // 字符串枚举不支持位运算，以下操作将失败
    added := country.Add(USA)
    fmt.Println("添加成功:", added)  // false - 字符串类型不支持添加操作
    
    removed := country.Remove(China)
    fmt.Println("移除成功:", removed)  // false - 字符串类型不支持移除操作
    
    // 获取枚举值和描述
    fmt.Println("枚举代码:", country.Code())           // CN
    fmt.Println("枚举描述:", country.Description())    // 中国
}
```

### 带数据的枚举示例

枚举值还可以携带额外的数据，适用于需要更复杂信息的场景：

```go
package main

import (
    "fmt"
    "github.com/kysion/base-library/utility/enum"
)

// 定义用户角色的额外数据结构
type RoleData struct {
    Level     int
    AllowList []string
}

func main() {
    // 创建带有额外数据的枚举
    adminRole := enum.NewWithData(
        1, 
        RoleData{Level: 10, AllowList: []string{"user", "role", "system"}},
        "管理员角色",
    )
    
    // 获取枚举的数据
    roleData := (*adminRole).Data()
    
    // 使用枚举数据
    fmt.Println("角色级别:", roleData.Level)
    fmt.Println("权限列表:", roleData.AllowList)
}
```
