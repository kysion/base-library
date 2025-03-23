# 基础模型 (base_model)

基础模型模块提供了数据模型的通用结构和方法，用于简化数据处理流程，实现一致的数据处理模式。

## 主要组件

基础模型模块包含以下主要组件：

### 基础模型结构 (BaseModel)

提供了通用的模型基础结构，包含常见的数据库字段和基本方法：

```go
type BaseModel struct {
    Id        uint      `json:"id"`         // 主键ID
    CreatedAt time.Time `json:"createdAt"`  // 创建时间
    UpdatedAt time.Time `json:"updatedAt"`  // 更新时间
    DeletedAt gtime.Time `json:"deletedAt"` // 删除时间（软删除）
    // 其他公共字段...
}
```

### 分页请求/响应结构

用于处理分页数据的请求和响应：

```go
// 分页请求参数
type PageReq struct {
    Page     int `json:"page" d:"1"`      // 当前页码，默认1
    PageSize int `json:"pageSize" d:"10"` // 每页记录数，默认10
}

// 分页响应结构
type PageRes struct {
    List     interface{} `json:"list"`               // 数据列表
    Page     int         `json:"page"`               // 当前页码
    PageSize int         `json:"pageSize"`           // 每页记录数
    Total    int         `json:"total"`              // 总记录数
    PageCount int        `json:"pageCount,omitempty"`// 总页数
}
```

### 响应结构

标准化的 API 响应结构：

```go
// 响应结构
type Response struct {
    Code    int         `json:"code"`    // 错误码，0表示成功
    Message string      `json:"message"` // 错误信息
    Data    interface{} `json:"data"`    // 响应数据
}
```

## 使用示例

### 基础模型的使用

通过嵌入 BaseModel 来扩展自定义模型：

```go
package models

import "github.com/kysion/base-library/base_model"

// 用户模型
type User struct {
    base_model.BaseModel             // 嵌入基础模型，自动包含 ID、时间等公共字段
    Username    string `json:"username"`    // 用户名
    Password    string `json:"-"`           // 密码，JSON 序列化时忽略
    Nickname    string `json:"nickname"`    // 昵称
    Email       string `json:"email"`       // 邮箱
    Phone       string `json:"phone"`       // 手机
    Status      int    `json:"status"`      // 状态
}
```

### 处理分页请求

使用分页请求和响应结构处理分页数据：

```go
import (
    "github.com/kysion/base-library/base_model"
    "github.com/gogf/gf/v2/frame/g"
)

func GetUserList(ctx context.Context, req *base_model.PageReq) (*base_model.PageRes, error) {
    // 创建分页响应
    res := &base_model.PageRes{
        Page:     req.Page,
        PageSize: req.PageSize,
    }
    
    // 查询数据
    var users []User
    result := g.DB().Model("users").
        Where("deleted_at IS NULL").
        Page(req.Page, req.PageSize).
        Scan(&users)
    if err := result.Err(); err != nil {
        return nil, err
    }
    
    // 获取总数
    count, err := g.DB().Model("users").Where("deleted_at IS NULL").Count()
    if err != nil {
        return nil, err
    }
    
    // 设置响应数据
    res.List = users
    res.Total = count
    res.PageCount = (count + req.PageSize - 1) / req.PageSize
    
    return res, nil
}
```

### 返回标准响应

使用标准响应结构返回 API 结果：

```go
import (
    "github.com/kysion/base-library/base_model"
    "github.com/kysion/base-library/base_consts"
)

func Success(data interface{}) *base_model.Response {
    return &base_model.Response{
        Code:    base_consts.ErrCodeSuccess,
        Message: "操作成功",
        Data:    data,
    }
}

func Error(code int, message string) *base_model.Response {
    return &base_model.Response{
        Code:    code,
        Message: message,
        Data:    nil,
    }
}

// 在接口中使用
func GetUser(r *ghttp.Request) {
    id := r.GetInt("id")
    user, err := service.GetUserById(r.Context(), id)
    if err != nil {
        r.Response.WriteJson(Error(base_consts.ErrCodeNotFound, "用户不存在"))
        return
    }
    r.Response.WriteJson(Success(user))
}
```

## 扩展模型

可以根据项目需要扩展基础模型，添加更多公共字段或方法：

```go
package project_model

import "github.com/kysion/base-library/base_model"

// 项目自定义基础模型
type ProjectBaseModel struct {
    base_model.BaseModel             // 嵌入基础模型
    CreatedBy  uint  `json:"createdBy"`  // 创建人ID
    UpdatedBy  uint  `json:"updatedBy"`  // 更新人ID
    DeletedBy  uint  `json:"deletedBy"`  // 删除人ID
    TenantId   uint  `json:"tenantId"`   // 租户ID
}

// 添加自定义方法
func (m *ProjectBaseModel) SetCreator(userId uint) {
    m.CreatedBy = userId
}

func (m *ProjectBaseModel) SetUpdater(userId uint) {
    m.UpdatedBy = userId
}
```

这样可以在项目中保持一致的数据模型结构，提高代码的可维护性和可重用性。
