package base_permission

import "fmt"

// IPermission 定义了权限相关的接口
type IPermission interface {
	GetId() int64
	GetParentId() int64
	GetName() string
	GetDescription() string
	GetIdentifier() string
	GetType() int
	GetMatchMode() int
	GetIsShow() int
	GetSort() int
	GetItems() []IPermission
	GetData() interface{}

	SetId(val int64) IPermission
	SetParentId(val int64) IPermission
	SetName(val string) IPermission
	SetDescription(val string) IPermission
	SetIdentifier(val string) IPermission
	SetType(val int) IPermission
	SetMatchMode(val int) IPermission
	SetIsShow(val int) IPermission
	SetSort(val int) IPermission
	SetItems(val []IPermission) IPermission
}

// Factory 需要在引用的业务项目中重新对 Factory 变量进行赋值
var (
	factory     func() IPermission
	initialized bool
)

// InitializePermissionFactory 初始化权限工厂，用于设置创建权限实例的方法
// 参数 permissionFactory 是一个函数，用于创建 IPPermission 实例
func InitializePermissionFactory(permissionFactory func() IPermission) {
	if permissionFactory != nil {
		factory = permissionFactory
		initialized = true
	}
}

// EnsureFactoryInitialized 确保工厂变量已被初始化
// 返回错误如果工厂未初始化
func EnsureFactoryInitialized() error {
	if !initialized {
		if factory == nil {
			return fmt.Errorf("InitializePermissionFactory must be initialized before use")
		}
		initialized = true
	}
	return nil
}

// New 构造权限信息
// 参数 id 权限的唯一标识符
// 参数 identifier 权限的标识符
// 参数 name 权限的名称
// 参数 description 权限的描述，可选
// 返回创建的权限实例和可能的错误
func New(id int64, identifier string, name string, description ...string) IPermission {
	err := EnsureFactoryInitialized()
	if err != nil {
		return nil
	}

	desc := ""
	if len(description) > 0 {
		desc = description[0]
	}

	perm := factory()
	if _, ok := perm.(IPermission); !ok {
		return nil
	}

	return perm.SetId(id).SetIdentifier(identifier).SetName(name).SetDescription(desc)
}

// NewInIdentifier 构造权限信息
// 参数 identifier 权限的标识符
// 参数 name 权限的名称
// 参数 description 权限的描述，可选
// 返回创建的权限实例和可能的错误
func NewInIdentifier(identifier string, name string, description ...string) IPermission {
	err := EnsureFactoryInitialized()
	if err != nil {
		return nil
	}

	desc := ""
	if len(description) > 0 {
		desc = description[0]
	}

	perm := factory()
	if _, ok := perm.(IPermission); !ok {
		return nil
	}

	return perm.SetIdentifier(identifier).SetName(name).SetDescription(desc)
}
