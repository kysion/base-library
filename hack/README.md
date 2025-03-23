# 工具脚本模块 (hack)

工具脚本模块提供了一系列实用的开发和运维工具脚本，用于简化代码生成、构建、测试、部署等常见任务，提高开发效率。

## 主要功能

工具脚本模块包含以下主要功能：

1. **代码生成**：自动生成常用代码结构和模板
2. **项目构建**：提供项目构建和打包脚本
3. **数据库工具**：数据库迁移、备份和恢复工具
4. **开发辅助**：开发环境配置和工具链
5. **部署脚本**：简化应用部署流程

## 可用脚本

### 1. 代码生成脚本

#### 模型生成器 (gen_model.sh)

根据数据库表结构自动生成模型代码：

```bash
# 使用方法
$ ./hack/gen_model.sh -d database_name -t table_name [-o output_dir] [-p package_name]

# 示例：生成用户表的模型
$ ./hack/gen_model.sh -d kysion_db -t users -p models
```

生成的模型示例：

```go
// 生成的模型代码
package models

import (
    "github.com/kysion/base-library/base_model"
    "time"
)

// Users 用户表
type Users struct {
    base_model.BaseModel
    Username  string    `json:"username" gorm:"column:username;type:varchar(50);not null;comment:'用户名'"`
    Password  string    `json:"-" gorm:"column:password;type:varchar(255);not null;comment:'密码'"`
    Email     string    `json:"email" gorm:"column:email;type:varchar(100);comment:'邮箱'"`
    Phone     string    `json:"phone" gorm:"column:phone;type:varchar(20);comment:'手机号'"`
    Status    int       `json:"status" gorm:"column:status;type:tinyint(1);not null;default:1;comment:'状态:0=禁用,1=启用'"`
}

// TableName 表名
func (Users) TableName() string {
    return "users"
}
```

#### API 生成器 (gen_api.sh)

根据模型生成 RESTful API 代码：

```bash
# 使用方法
$ ./hack/gen_api.sh -m model_name [-o output_dir] [-p package_name]

# 示例：生成用户 API
$ ./hack/gen_api.sh -m Users -p api/user
```

生成的 API 控制器示例：

```go
// 生成的控制器代码
package user

import (
    "github.com/gogf/gf/v2/net/ghttp"
    "github.com/kysion/base-library/base_model"
    "my-project/models"
    "my-project/service"
)

// Controller 用户控制器
type Controller struct{}

// List 获取用户列表
func (c *Controller) List(r *ghttp.Request) {
    var req base_model.PageReq
    if err := r.Parse(&req); err != nil {
        r.Response.WriteJson(base_model.Error(1, err.Error()))
        return
    }
    
    result, err := service.User.GetList(r.Context(), &req)
    if err != nil {
        r.Response.WriteJson(base_model.Error(1, err.Error()))
        return
    }
    
    r.Response.WriteJson(base_model.Success(result))
}

// Get 获取用户详情
func (c *Controller) Get(r *ghttp.Request) {
    id := r.GetInt("id")
    if id <= 0 {
        r.Response.WriteJson(base_model.Error(1, "无效的用户ID"))
        return
    }
    
    user, err := service.User.GetById(r.Context(), id)
    if err != nil {
        r.Response.WriteJson(base_model.Error(1, err.Error()))
        return
    }
    
    r.Response.WriteJson(base_model.Success(user))
}

// Create 创建用户
func (c *Controller) Create(r *ghttp.Request) {
    var user models.Users
    if err := r.Parse(&user); err != nil {
        r.Response.WriteJson(base_model.Error(1, err.Error()))
        return
    }
    
    if err := service.User.Create(r.Context(), &user); err != nil {
        r.Response.WriteJson(base_model.Error(1, err.Error()))
        return
    }
    
    r.Response.WriteJson(base_model.Success(user))
}

// Update 更新用户
func (c *Controller) Update(r *ghttp.Request) {
    var user models.Users
    if err := r.Parse(&user); err != nil {
        r.Response.WriteJson(base_model.Error(1, err.Error()))
        return
    }
    
    if user.Id <= 0 {
        r.Response.WriteJson(base_model.Error(1, "无效的用户ID"))
        return
    }
    
    if err := service.User.Update(r.Context(), &user); err != nil {
        r.Response.WriteJson(base_model.Error(1, err.Error()))
        return
    }
    
    r.Response.WriteJson(base_model.Success(nil))
}

// Delete 删除用户
func (c *Controller) Delete(r *ghttp.Request) {
    id := r.GetInt("id")
    if id <= 0 {
        r.Response.WriteJson(base_model.Error(1, "无效的用户ID"))
        return
    }
    
    if err := service.User.Delete(r.Context(), id); err != nil {
        r.Response.WriteJson(base_model.Error(1, err.Error()))
        return
    }
    
    r.Response.WriteJson(base_model.Success(nil))
}
```

### 2. 构建和部署脚本

#### 项目构建脚本 (build.sh)

编译并打包应用程序：

```bash
# 使用方法
$ ./hack/build.sh [-e env] [-v version] [-o output_dir]

# 示例：构建生产环境版本
$ ./hack/build.sh -e prod -v 1.0.0
```

#### 数据库迁移脚本 (migrate.sh)

管理数据库结构变更：

```bash
# 使用方法
$ ./hack/migrate.sh [-c command] [-d database] [-f file]

# 示例：创建迁移文件
$ ./hack/migrate.sh -c create -f create_users_table

# 示例：应用迁移
$ ./hack/migrate.sh -c up -d kysion_db

# 示例：回滚迁移
$ ./hack/migrate.sh -c down -d kysion_db
```

#### Docker 部署脚本 (docker_deploy.sh)

使用 Docker 部署应用：

```bash
# 使用方法
$ ./hack/docker_deploy.sh [-t tag] [-e env]

# 示例：部署到测试环境
$ ./hack/docker_deploy.sh -t v1.0.0 -e test
```

### 3. 开发辅助工具

#### 开发环境配置脚本 (setup_dev.sh)

快速配置开发环境：

```bash
# 使用方法
$ ./hack/setup_dev.sh

# 该脚本会：
# 1. 安装必要的依赖和工具
# 2. 配置环境变量
# 3. 创建本地开发数据库
# 4. 应用初始迁移
```

#### 代码检查脚本 (lint.sh)

检查代码质量和格式：

```bash
# 使用方法
$ ./hack/lint.sh [dir]

# 示例：检查所有代码
$ ./hack/lint.sh

# 示例：检查特定目录
$ ./hack/lint.sh ./api
```

#### 测试运行脚本 (test.sh)

运行单元测试和集成测试：

```bash
# 使用方法
$ ./hack/test.sh [-t test_type] [-c coverage]

# 示例：运行单元测试
$ ./hack/test.sh -t unit

# 示例：运行集成测试并生成覆盖率报告
$ ./hack/test.sh -t integration -c true
```

## 扩展脚本

你可以根据项目需求扩展和自定义脚本：

### 创建自定义脚本

1. 在 `hack` 目录下创建新的脚本文件
2. 添加脚本头部信息和使用说明
3. 实现脚本功能
4. 添加错误处理和日志输出

示例：

```bash
#!/bin/bash

# 脚本：custom_script.sh
# 描述：这是一个自定义脚本示例
# 用法：./hack/custom_script.sh [options]

# 显示帮助信息
function show_help() {
    echo "用法: $0 [options]"
    echo "选项:"
    echo "  -h, --help     显示帮助信息"
    echo "  -o, --option   自定义选项"
    exit 0
}

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    key="$1"
    case $key in
        -h|--help)
            show_help
            ;;
        -o|--option)
            OPTION="$2"
            shift
            ;;
        *)
            echo "未知选项: $1"
            show_help
            ;;
    esac
    shift
done

# 脚本主逻辑
echo "执行自定义脚本..."
# 实现你的功能
echo "完成！"
```

## 最佳实践

1. **脚本文档**：为每个脚本提供清晰的文档和使用示例
2. **参数校验**：始终验证脚本输入参数，提供有意义的错误消息
3. **错误处理**：妥善处理脚本执行过程中的错误，避免中断
4. **日志输出**：使用统一的日志格式，方便问题排查
5. **保持简单**：每个脚本专注于单一功能，避免过度复杂
6. **版本控制**：将脚本纳入版本控制，记录变更历史
7. **环境兼容**：确保脚本在不同环境中表现一致
