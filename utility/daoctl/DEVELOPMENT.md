# daoctl 开发指南

本文档提供了 `daoctl` 工具包的详细开发指南，包括代码结构、开发规范和贡献指导。

## 代码结构

`daoctl` 工具包的代码结构如下：

```
daoctl/
├── cache.go           # 缓存相关功能实现
├── dao.go             # DAO 核心功能和配置
├── dao_interface/     # 接口定义
│   └── interface.go   # IDao 和 TIDao 接口
├── delete.go          # 删除操作实现
├── index.go           # 主要查询功能
├── insert.go          # 插入操作实现
├── internal/          # 内部工具函数
│   └── ...
├── save.go            # 保存操作实现
├── scan.go            # 扫描结果实现
└── update.go          # 更新操作实现
```

## 核心模块说明

### dao.go

包含 DAO 的核心配置和功能：

- `DaoConfig` 结构体的初始化和管理
- ORM 缓存的开启与关闭
- 扩展模型映射的管理
- 上下文相关的操作

### index.go

包含主要的查询功能：

- `GetById`：根据 ID 获取单个记录
- `Find`：根据条件查找记录
- `GetAll`：获取全部记录
- `Query`：复杂查询操作

### dao_interface/interface.go

定义了 DAO 操作的核心接口：

- `IDao`：基础 DAO 接口
- `TIDao`：支持泛型的 DAO 接口
- `DaoConfig`：DAO 配置结构

### internal/

内部工具函数，不对外暴露：

- 查询条件构建
- 排序处理
- 参数验证

## 开发规范

### 命名规范

- 函数名：使用驼峰命名法，如 `GetById`、`FindOne`
- 变量名：使用驼峰命名法，如 `dataModel`、`cacheOption`
- 常量：使用下划线分隔的大写字母，如 `CONTEXT_MODEL_TABLE_KEY`
- 文件名：使用小写字母，如 `dao.go`、`index.go`

### 注释规范

每个导出的函数、类型和常量必须有注释，遵循以下格式：

```go
// FunctionName 简短描述。
//
// 详细描述，可以跨多行。
// 参数:
// - param1: 参数1的描述。
// - param2: 参数2的描述。
//
// 返回值:
// - 返回值1的描述。
// - 返回值2的描述。
func FunctionName(param1 string, param2 int) (string, error) {
    // 实现...
}
```

### 错误处理

- 所有可能返回错误的函数必须正确处理错误
- 错误信息应当清晰、具体，便于定位问题
- 使用 `fmt.Errorf` 在转发错误时添加上下文信息

```go
if err := doSomething(); err != nil {
    return nil, fmt.Errorf("failed to do something: %w", err)
}
```

### 代码格式

使用 `gofmt` 或 `goimports` 进行代码格式化：

```bash
gofmt -w .
goimports -w .
```

## 贡献指南

### 准备工作

1. Fork 项目仓库
2. 克隆你的 Fork 到本地
3. 创建新的分支：`git checkout -b feature/your-feature-name`

### 开发流程

1. 实现新功能或修复 bug
2. 添加或更新测试用例
3. 确保所有测试通过
4. 更新相关文档
5. 提交代码并创建 Pull Request

### 提交规范

- 使用清晰、具体的提交信息
- 每个提交专注于单一变更
- 提交信息格式：`类型: 简短描述`

类型可以是：

- `feat`：新功能
- `fix`：Bug 修复
- `docs`：文档更新
- `style`：代码格式变更
- `refactor`：代码重构
- `perf`：性能优化
- `test`：测试相关
- `chore`：构建过程或辅助工具变动

示例：`feat: 添加缓存自动失效机制`

### 测试要求

- 新功能必须有对应的单元测试
- 修复 bug 必须有对应的回归测试
- 测试覆盖率应保持在 70% 以上

## 扩展开发指南

### 添加新的查询方法

1. 在适当的文件中添加新函数，如查询相关功能在 `index.go`
2. 遵循现有的泛型模式
3. 添加详细注释
4. 编写单元测试

示例：

```go
// FindByField 根据指定字段查找记录。
//
// 参数:
// - model: 数据库模型。
// - field: 查询字段名。
// - value: 查询值。
//
// 返回值:
// - 返回查询到的记录集合和可能的错误。
func FindByField[T any](model *gdb.Model, field string, value interface{}) (*base_model.CollectRes[T], error) {
    model = ExecExWhere(model)
    return Find[T](model, nil, base_model.FilterInfo{
        Field:    field,
        Operator: "=",
        Value:    value,
    })
}
```

### 添加新的缓存策略

1. 在 `cache.go` 中定义新的缓存函数
2. 确保缓存策略可配置
3. 添加对应的单元测试

示例：

```go
// MakeCustomCache 创建自定义缓存选项。
//
// 参数:
// - key: 缓存键名。
// - duration: 缓存有效期。
//
// 返回值:
// - 返回自定义缓存选项。
func MakeCustomCache(key string, duration time.Duration) *gdb.CacheOption {
    return &gdb.CacheOption{
        Duration: duration,
        Name:     "custom:" + key,
        Force:    false,
    }
}
```

### 添加扩展模型处理器

1. 定义适用于特定表的处理函数
2. 使用 `MakeExtModelMap` 注册处理函数
3. 确保处理函数可以通过上下文禁用

示例：

```go
// 注册用户表的软删除过滤器
daoctl.MakeExtModelMap("user", map[string]func(model *gdb.Model, conf *dao_interface.DaoConfig, data ...interface{}) *gdb.Model{
    "filter_deleted": func(model *gdb.Model, conf *dao_interface.DaoConfig, data ...interface{}) *gdb.Model {
        return model.Where("deleted_at", nil)
    },
})
```

## 性能优化建议

1. **减少数据库查询次数**
   - 使用联合查询代替多次单表查询
   - 合理使用 `IN` 条件一次获取多条记录

2. **合理使用缓存**
   - 频繁查询的数据优先使用缓存
   - 大量数据考虑分段缓存

3. **优化查询条件**
   - 建立适当的索引
   - 减少不必要的字段查询，使用列投影

4. **减少内存分配**
   - 预分配切片容量
   - 重用对象而非频繁创建

5. **并发控制**
   - 使用连接池控制并发连接数
   - 耗时操作考虑异步处理

## 常见问题与解决方案

### 缓存不生效

1. 检查是否正确启用了缓存：`EnableOrmCache()`
2. 确认表名不在忽略缓存的配置中
3. 验证 `CacheOption` 是否正确传递

### 查询结果不符合预期

1. 检查是否有扩展模型条件影响了查询
2. 使用 `IgnoreExtModel` 方法禁用特定的扩展条件
3. 检查泛型类型是否与表结构匹配

### 性能问题

1. 检查是否有过多的数据库查询
2. 确认是否正确使用了缓存
3. 考虑优化查询条件或增加索引

## 参考资料

- [GoFrame 官方文档](https://goframe.org/pages/viewpage.action?pageId=1114119)
- [Go 泛型指南](https://go.dev/doc/tutorial/generics)
- [数据库索引设计指南](https://use-the-index-luke.com/)
- [Go 代码规范](https://github.com/golang/go/wiki/CodeReviewComments)
