package base_tree

// TestTree 结构实现了 Tree 接口，用于表示具有层次结构的树形数据。
// 它包含节点信息、父节点ID、节点名称以及子节点列表。
type TestTree struct {
	Id       int64       `json:"id"             dc:"ID" v:"integer"`                   // 节点ID，唯一标识一个节点。
	ParentId int64       `json:"parentId"       dc:"父级ID" v:"min:0#必须是正整数"`            // 父节点ID，标识该节点的父节点。
	Name     string      `json:"name"           dc:"名称" v:"max-length:64#仅支持最大字符长度64"` // 节点名称，最大长度为64个字符。
	Children []*TestTree `json:"children"       dc:"子树"`                               // 子节点列表，表示该节点的所有直接子节点。
}

// IsParentChildEqual 检查给定的父节点和子节点ID是否匹配。
// 参数:
//
//	father: 指向一个父节点的指针。
//	childId: 需要检查的子节点ID。
//
// 返回值:
//
//	如果父节点和子节点ID匹配，则返回true；否则返回false。
func (d *TestTree) IsParentChildEqual(father *TestTree, childId int64) bool {
	if father == nil {
		return false
	}
	return father.Id == childId
}

// AssignChildren 将给定的子节点列表分配给指定的父节点。
// 参数:
//
//	father: 指向一个父节点的指针，将子节点列表分配给它。
//	children: 指向子节点列表的指针，这些子节点将被分配给父节点。
func (d *TestTree) AssignChildren(father *TestTree, children []*TestTree) {
	if father != nil {
		father.Children = children
	}
}

// IsTopLevel 检查指定的父节点是否为顶级节点。
// 参数:
//
//	father: 指向一个父节点的指针，检查它是否为顶级节点。
//
// 返回值:
//
//	如果父节点的父ID为0，则返回true，表示它是顶级节点；否则返回false。
func (d *TestTree) IsTopLevel(father *TestTree) bool {
	if father == nil {
		return false
	}
	return father.ParentId == 0
}
