// ================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package share_service

{Imports}
`

const TemplateGenServiceContentInterface = `
{InterfaceName} interface {
	{FuncDefinition}
}
`

const TemplateGenServiceContentVariable = `
local{StructName} {InterfaceName}
`

const TemplateGenServiceContentRegister = `
func {StructName}() {InterfaceName} {
	if local{StructName} == nil {
		panic("implement not found for interface {InterfaceName}, forgot register?")
	}
	return local{StructName}
}

func Register{StructName}(i {InterfaceName}) {
	local{StructName} = i
}