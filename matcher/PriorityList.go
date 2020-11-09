package matcher

import (
	"errors"
	"fmt"
)

type ISortable interface {
	Less(sortable ISortable) bool
	Equal(sortable ISortable) bool
}

type IEnumerator interface {
	OnReceived(index int, current ISortable) bool
	OnAbolished(index int, current ISortable)
}

type PriorityList struct {
	mSlice []ISortable
}

func (pl *PriorityList) insertAtIndex(itemToInsert ISortable, index int) {
	afterItems := append([]ISortable{}, pl.mSlice[index:]...)
	pl.mSlice = append(pl.mSlice[0:index], itemToInsert)
	pl.mSlice = append(pl.mSlice, afterItems...)
}

// 倒序添加（从大到小排序）
func (pl *PriorityList) AddZ(itemToAdd ISortable) error {
	if pl.mSlice == nil {
		pl.insertAtIndex(itemToAdd, 0)
		return nil
	}
	var index int = 0
	for _, item := range pl.mSlice {
		if itemToAdd.Equal(item) { //重复添加
			var errMsg = fmt.Sprintf("item:%v has a same priority", item)
			return errors.New(errMsg)
		}
		if !itemToAdd.Less(item) {
			break
		}

		index++
	}
	pl.insertAtIndex(itemToAdd, index)
	return nil
}

// 按照优先级大小倒序遍历
func (pl PriorityList) EnumerateZ(callback func(index int, item ISortable) bool) {
	var index int = 0
	for _, item := range pl.mSlice { //已经倒序排序，直接遍历
		var needToAbolish = callback(index, item)
		if needToAbolish {
			break
		}
		index++
	}
}
