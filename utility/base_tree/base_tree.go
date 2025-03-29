package base_tree

// Filter 过滤并构建树List
// T 为任意类型
// arr 为待过滤的数组
// f 为过滤条件函数，对每个元素进行判断，返回bool值
// 返回值为过滤后的新列表
func Filter[T any](arr []T, f func(item T) bool) (list []T) {
	for _, el := range arr {
		// 当满足过滤条件时，将元素添加到列表中
		if f(el) {
			list = append(list, el)
		}
	}
	return list
}

// Tree 接口定义了树操作的抽象
type Tree[T any] interface {
	// IsParentChildEqual 用于判断child.ParentId是否等于father.Id
	// father 为父节点
	// child 为子节点
	// 返回值为bool，表示是否相等
	IsParentChildEqual(father T, child T) bool
	// AssignChildren 将子树List设置到父树
	// father 为父节点
	// branchArr 为子树列表
	AssignChildren(father T, branchArr []T)
	// IsRoot 用于判断father.ID是否等于指定ParentID，可作为递归使用场景的退出条件
	// father 为父节点
	// 返回值为bool，表示是否满足退出条件
	IsRoot(father T) bool
	// MakeSubNodeSort 对子节点进行排序
	MakeSubNodeSort()
}

// ToTree 将列表转换为树结构
// list 为待转换的列表
// fun 为实现Tree接口的具体类型
// 返回值为树结构列表
func ToTree[T any](list []T, fun Tree[T]) []T {
	// 外层递归，返回Tree树结构 （父 + 子）
	return Filter(list, func(father T) bool {
		// 内层递归，用于构建子树，
		branchArr := Filter(list, func(childId T) bool {
			// 内层递归构建树List的条件：fun.GeteIsQual(father, childId)，当Child的FacterId == FacterId的时候
			return fun.IsParentChildEqual(father, childId)
		})

		// 如果子树Arr不为空， 那么父树father.Children = branchArr设置上
		if len(branchArr) > 0 {
			fun.AssignChildren(father, branchArr)
		}

		// 外层递归退出条件是：father.ID = 0
		return fun.IsRoot(father)
	})
}

