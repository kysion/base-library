package base_permission

import "github.com/yitter/idgenerator-go/idgen"

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

var PFactory func() IPermission

var Factory = PFactory

// New 构造权限信息
func New(id int64, identifier string, name string, description ...string) IPermission {
	var desc string
	if len(description) > 0 {
		desc = description[0]
	}

	return Factory().SetId(id).SetIdentifier(identifier).SetName(name).SetDescription(desc)
}

// NewInIdentifier 构造权限信息
func NewInIdentifier(identifier string, name string, description ...string) IPermission {
	var desc string

	if len(description) > 0 {
		desc = description[0]
	}

	return Factory().SetId(idgen.NextId()).SetIdentifier(identifier).SetName(name).SetDescription(desc)
}
