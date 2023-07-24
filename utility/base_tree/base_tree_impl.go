package base_tree

/*
	Tree接口的实现示例代码
*/

type testTree struct {
	Id       int64  `json:"id"             dc:"ID" v:"integer"`
	ParentId int64  `json:"parentId"       dc:"父级ID" v:"min:0#必须是正整数，该属性创建后不支持修改"`
	Name     string `json:"name"           dc:"名称" v:"max-length:64#仅支持最大字符长度64"`
}

type testPermissionTree struct {
	*testTree
	Children []*testPermissionTree `json:"children"       dc:"子树"`
}

func (d *testPermissionTree) GetIsEqual(father *testPermissionTree, childId *testPermissionTree) bool {
	return father.Id == childId.ParentId
}
func (d *testPermissionTree) SetChild(father *testPermissionTree, branchArr []*testPermissionTree) {
	father.Children = branchArr
}
func (d *testPermissionTree) RetFather(father *testPermissionTree) bool {
	// 顶级的ParentId这块可以看一下保存的时候ParentId 默认值是多少
	return father.ParentId == 0
}
