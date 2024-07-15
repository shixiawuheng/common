package common

import (
	"container/list"
)

type RedisList struct {
	l *list.List
}

func New() *RedisList {
	return &RedisList{l: list.New()}
}

func (rl *RedisList) LPush(values ...interface{}) int {
	for _, val := range values {
		rl.l.PushFront(val)
	}
	return rl.l.Len()
}

func (rl *RedisList) RPush(values ...interface{}) int {
	for _, val := range values {
		rl.l.PushBack(val)
	}
	return rl.l.Len()
}

func (rl *RedisList) LPop() interface{} {
	if rl.l.Len() == 0 {
		return nil
	}
	return rl.l.Remove(rl.l.Front())
}

func (rl *RedisList) RPop() interface{} {
	if rl.l.Len() == 0 {
		return nil
	}
	return rl.l.Remove(rl.l.Back())
}

func (rl *RedisList) LLen() int {
	return rl.l.Len()
}

func (rl *RedisList) LIndex(index int) interface{} {
	// 转为正数索引
	if index < 0 {
		index = rl.l.Len() + index
	}
	// 判断索引是否越界
	if index < 0 || index >= rl.l.Len() {
		return nil
	}
	// 遍历list取值
	i := 0
	for e := rl.l.Front(); e != nil; e = e.Next() {
		if i == index {
			return e.Value
		}
		i++
	}
	return nil
}

func (rl *RedisList) LRem(value interface{}, count int) int {
	var removedCount int
	if count > 0 {
		for e, i := rl.l.Front(), 0; e != nil && i < count; {
			if e.Value == value {
				next := e.Next()
				rl.l.Remove(e)
				e = next
				removedCount++
			} else {
				e = e.Next()
			}
		}
	} else if count < 0 {
		for e, i := rl.l.Back(), 0; e != nil && i < -count; {
			if e.Value == value {
				prev := e.Prev()
				rl.l.Remove(e)
				e = prev
				removedCount++
			} else {
				e = e.Prev()
			}
		}
	} else {
		for e := rl.l.Front(); e != nil; {
			if e.Value == value {
				next := e.Next()
				rl.l.Remove(e)
				e = next
				removedCount++
			} else {
				e = e.Next()
			}
		}
	}
	return removedCount
}

func (rl *RedisList) LRange(start, stop int) []interface{} {
	// 转为正数索引
	if start < 0 {
		start = rl.l.Len() + start
	}
	if stop < 0 {
		stop = rl.l.Len() + stop
	}
	// 判断索引是否越界
	if start < 0 {
		start = 0
	}
	if stop >= rl.l.Len() {
		stop = rl.l.Len() - 1
	}
	if start > stop {
		return nil
	}
	// 遍历list取值
	res := make([]interface{}, stop-start+1)
	i := 0
	for e := rl.l.Front(); e != nil && i <= stop; e = e.Next() {
		if i >= start {
			res[i-start] = e.Value
		}
		i++
	}
	return res
}
