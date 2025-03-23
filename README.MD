# Kysion Base Library 基础库

Kysion Base Library 是一个综合性的 Go 语言工具库，提供了可复用的组件和工具函数，简化应用开发流程，提高代码质量。本库基于 Go 1.18+ 实现，充分利用了泛型特性，确保类型安全。

[![Go Report Card](https://goreportcard.com/badge/github.com/kysion/base-library)](https://goreportcard.com/report/github.com/kysion/base-library)
[![GoDoc](https://godoc.org/github.com/kysion/base-library?status.svg)](https://godoc.org/github.com/kysion/base-library)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## 目录

- [功能概览](#功能概览)
- [安装](#安装)
- [使用指南](#使用指南)
- [模块详解](#模块详解)
- [性能表现](#性能表现)
- [CI/CD 自动化](#cicd-自动化)
- [版本管理](#版本管理)
- [版本兼容性](#版本兼容性)
- [贡献指南](#贡献指南)
- [许可证](#许可证)

## 功能概览

基础库包含以下主要模块：

- **基础常量 (base_consts)** - 提供常用的系统常量定义
- **基础模型 (base_model)** - 提供数据模型的基础结构和方法
- **基础钩子 (base_hook)** - 提供插件式扩展的钩子系统
- **工具集 (utility)** - 提供各类实用工具函数和结构
  - **枚举工具 (enum)** - 提供类型安全的枚举实现
  - **DAO 工具 (daoctl)** - 提供数据访问对象控制工具
  - **随机工具 (base_random)** - 提供高效安全的随机生成功能

## 安装

通过 Go 模块引入库:

```bash
go get -u github.com/kysion/base-library@latest
```

## 使用指南

以下是一些常用模块的基本使用示例。

### 基础常量 (base_consts)

```go
import "github.com/kysion/base-library/base_consts"

func main() {
    // 使用预定义常量
    status := base_consts.StatusActive
}
```

### 基础模型 (base_model)

```go
import "github.com/kysion/base-library/base_model"

// 使用基础模型创建自定义模型
type User struct {
    base_model.BaseModel
    Name string
    Age  int
}
```

### 基础钩子 (base_hook)

```go
import "github.com/kysion/base-library/base_hook"

func main() {
    // 创建钩子管理器
    hooks := base_hook.NewManager()
    
    // 注册钩子处理函数
    hooks.Register("before_save", func(data interface{}) error {
        // 处理逻辑
        return nil
    })
    
    // 触发钩子
    hooks.Trigger("before_save", userData)
}
```

### 工具集 (utility)

#### 枚举工具 (enum)

```go
import "github.com/kysion/base-library/utility/enum"

// 定义权限枚举
type Permission interface {
    enum.IEnumCode[int]
}

// 创建枚举值
var (
    ReadPerm = enum.New[Permission](1, "读取权限")
    WritePerm = enum.New[Permission](2, "写入权限")
)

func main() {
    // 创建组合权限
    userPerm := enum.New[Permission](0, "用户权限")
    userPerm.Add(ReadPerm, WritePerm)
    
    // 检查权限
    hasReadPerm := userPerm.Has(ReadPerm) // true
}
```

#### DAO工具 (daoctl)

```go
import (
    "github.com/kysion/base-library/base_model"
    "github.com/kysion/base-library/utility/daoctl"
    "github.com/gogf/gf/v2/database/gdb"
)

func ExampleUseDao() {
    // 假设已有数据库模型
    model := g.DB().Model("users")
    
    // 构建查询条件
    filter := []base_model.FilterInfo{
        {Field: "name", Where: "like", Value: "%张%"},
        {Field: "age", Where: ">", Value: 18},
    }
    
    // 构建排序条件
    orderBy := []base_model.OrderBy{
        {Field: "created_at", Sort: "desc"},
    }
    
    // 执行查询
    result, err := daoctl.Find[YourEntityType](model, orderBy, filter...)
    if err != nil {
        // 处理错误
    }
    
    // 使用查询结果
    for _, item := range result.Records {
        // 处理每条记录
    }
}
```

#### 随机工具 (base_random)

```go
import "github.com/kysion/base-library/utility/base_random"

func main() {
    // 生成随机字符串
    randomStr := base_random.GenerateRandomString(10) // 生成10位随机字符串
    
    // 生成随机数字字符串
    randomNum := base_random.GenerateNumberString(6) // 生成6位随机数字字符串
}
```

## 模块详解

- [**基础常量 (base_consts)**](base_consts/README.md) - 提供常用的系统常量定义
- [**基础模型 (base_model)**](base_model/README.md) - 提供数据模型的基础结构和方法
- [**基础钩子 (base_hook)**](base_hook/README.md) - 提供插件式扩展的钩子系统
- [**工具集 (utility)**](utility/README.md) - 提供各类实用工具函数和结构
  - [**枚举工具 (enum)**](utility/enum/README.md) - 提供类型安全的枚举实现
  - [**DAO工具 (daoctl)**](utility/daoctl/README.md) - 提供数据访问对象控制工具
  - [**随机工具 (base_random)**](utility/base_random/README.md) - 提供随机生成功能

## 性能表现

Kysion Base Library 注重性能优化，各组件经过严格的基准测试验证。以下是枚举工具包的性能测试结果摘要：

| 操作类型                | 操作耗时 (ns)   | 内存分配      |
|------------------------|----------------|--------------|
| 基本枚举访问操作        | < 0.2 ns       | 零内存分配    |
| 单枚举位运算操作        | < 0.8 ns       | 零内存分配    |
| 多枚举位运算操作        | < 1.5 ns       | 零内存分配    |
| ToMap转换操作          | ~ 25.7 ns      | 16B/操作     |
| 并行混合操作            | ~ 6.9 ns       | 16B/操作     |

上述测试在 Intel i9-14900K 处理器上运行，完整的性能测试报告请查看 [utility/enum/benchmark/README.md](utility/enum/benchmark/README.md)。

### 性能特点

- 核心操作高效：基础操作执行速度极快，每秒可执行数十亿次操作
- 内存友好：大多数操作零内存分配，有效减少GC压力
- 并发支持：在高并发环境中表现稳定
- 类型安全：通过泛型和类型断言机制确保类型安全，不牺牲性能

## CI/CD 自动化

本项目配置了完整的 CI/CD 自动化流程，确保代码质量和版本管理的标准化。

### 持续集成 (CI)

当代码被推送到仓库或创建 Pull Request 时，会自动运行以下检查：

- **代码规范检查**：使用 golangci-lint 检查代码质量
- **单元测试**：运行所有测试并生成覆盖率报告
- **安全扫描**：使用 gosec 进行代码安全漏洞扫描

### 持续部署 (CD)

当代码合并到主分支时，会自动触发版本发布流程：

1. 检出代码
2. 运行测试验证代码质量
3. 自动增加版本号（默认为补丁版本+1）
4. 生成变更日志
5. 创建 GitHub Release
6. 更新依赖并提交更改

## 版本管理

本项目使用[语义化版本](https://semver.org/lang/zh-CN/)进行版本管理：

- **主版本号 (MAJOR)**：进行不兼容的 API 修改时递增
- **次版本号 (MINOR)**：新增向下兼容的功能时递增
- **修订号 (PATCH)**：进行向下兼容的问题修正时递增

### 手动触发版本发布

可以通过 GitHub Actions 手动触发版本发布，并指定版本更新类型：

1. 在 GitHub 仓库页面选择 "Actions" 标签
2. 选择 "Release Go Module" 工作流
3. 点击 "Run workflow" 按钮
4. 选择所需的版本更新类型（patch/minor/major）
5. 点击 "Run workflow" 开始发布流程

## 版本兼容性

- Go 1.24+ (需要较新的泛型支持)
- 依赖 [GoFrame v2](https://github.com/gogf/gf) 框架

## 贡献指南

我们欢迎各种形式的贡献，包括但不限于功能增强、bug修复、文档改进等。

贡献流程：

1. Fork 项目仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送分支 (`git push origin feature/amazing-feature`)
5. 提交 Pull Request

请确保您的代码符合以下要求：

- 通过所有测试
- 遵循项目的代码规范
- 包含适当的文档和注释
- 确保类型安全，特别是在使用泛型时

## 许可证

本项目采用 MIT 许可证 - 详情请查看 [LICENSE](LICENSE) 文件。
