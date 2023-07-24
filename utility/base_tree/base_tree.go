package base_tree

// Filter 过滤并构建树List
func Filter[T any](arr []*T, f func(item *T) bool) (list []*T) {
	for _, el := range arr {
		// 当Child的FacterId == FacterId
		if f(el) {
			list = append(list, el)
		}
	}
	return list
}

type Tree[T any] interface {
	// GetIsEqual 用于判断child.ParentId是否等于facter.Id父ID
	GetIsEqual(father *T, child *T) bool
	// SetChild 将子树List设置到父树
	SetChild(father *T, branchArr []*T)
	// RetFather 用于判断father.ID是否等于指定ParentID，可作为递归使用场景的退出条件
	RetFather(father *T) bool
}

func ToTree[T any](list []*T, fun Tree[T]) []*T {

	// 外层递归，返回Tree树结构 （父 + 子）
	return Filter(list, func(father *T) bool {
		// 内层递归，用于构建子树，
		branchArr := Filter(list, func(childId *T) bool {
			// 内层递归构建树List的条件：fun.GeteIsQual(father, childId)，当Child的FacterId == FacterId的时候
			return fun.GetIsEqual(father, childId)
		})

		// 如果子树Arr不为空， 那么父树father.Children = branchArr设置上
		if len(branchArr) > 0 {
			fun.SetChild(father, branchArr)
		}

		// 外层递归退出条件是：father.ID = 0
		return fun.RetFather(father)
	})
}
