package sys_enum_hook

type hook struct {
	BusinessType businessType
}

var Hook = hook{
	BusinessType: BusinessType,
}
