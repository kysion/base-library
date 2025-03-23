# 枚举工具包 (Enum Utility)

这个工具包提供了一个灵活的枚举类型实现，支持整型、字符串型等多种类型的枚举值，以及枚举之间的位运算操作。

## 主要特性

- 支持多种数据类型的枚举值（整型、字符串型等）
- 支持枚举值之间的位运算（检查、添加、移除）
- 支持带有附加数据的枚举类型
- 高性能实现，基本操作无内存分配

## 目录结构

- `enum.go` - 枚举工具包的主要实现
- `benchmark/` - 性能测试相关文件
  - `enum_test.go` - 性能测试代码
  - `run.sh` - 运行性能测试的脚本
  - `README.md` - 性能测试文档和结果

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

## 性能测试

查看 [benchmark/README.md](benchmark/README.md) 了解性能测试结果和分析。
