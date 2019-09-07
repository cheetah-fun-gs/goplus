// Package sort 排序方法, 不推荐使用
package sort

import "sort"

// Any 可排序对象
type Any struct {
	Weight int64
	Object interface{}
}

// AnySortList 排序列表
type AnySortList []*Any

// Len Len
func (anySortList AnySortList) Len() int { return len(anySortList) }

// Swap Swap
func (anySortList AnySortList) Swap(i, j int) {
	anySortList[i], anySortList[j] = anySortList[j], anySortList[i]
}

// Less Less
func (anySortList AnySortList) Less(i, j int) bool {
	return anySortList[i].Weight < anySortList[j].Weight
}

// AnySortReverseList 排序列表 逆序
type AnySortReverseList []*Any

// Len Len
func (anySortList AnySortReverseList) Len() int { return len(anySortList) }

// Swap Swap
func (anySortList AnySortReverseList) Swap(i, j int) {
	anySortList[i], anySortList[j] = anySortList[j], anySortList[i]
}

// Less Less
func (anySortList AnySortReverseList) Less(i, j int) bool {
	return anySortList[i].Weight > anySortList[j].Weight
}

// AnyList 排序
func AnyList(anyList []*Any, reverse bool) {
	var sortList sort.Interface
	if !reverse {
		sortList = AnySortList(anyList)
	} else {
		sortList = AnySortReverseList(anyList)
	}
	sort.Sort(sortList)
	return
}
