package filer

import (
	"math"
	"sync"
)

// IntervalValue 是所有区间数据必须实现的接口，
// 提供设置区间起止位置的方法以及深拷贝方法。
type IntervalValue interface {
	SetStartStop(start, stop int64)
	Clone() IntervalValue
}

// Interval 表示一个区间数据节点。
// T 是实现了 IntervalValue 接口的泛型类型。
type Interval[T IntervalValue] struct {
	StartOffset int64        // 区间开始偏移量
	StopOffset  int64        // 区间结束偏移量
	TsNs        int64        // 时间戳，纳秒，用于版本比较（新旧）
	Value       T            // 区间值，必须实现 IntervalValue 接口
	Prev        *Interval[T] // 前一个区间节点
	Next        *Interval[T] // 后一个区间节点
}

// 计算区间大小（结束-起始）
func (interval *Interval[T]) Size() int64 {
	return interval.StopOffset - interval.StartOffset
}

// IntervalList 是一个有序的区间链表结构，线程安全。
// 用于表示某个文件的一系列非重叠写入区间。
type IntervalList[T IntervalValue] struct {
	head *Interval[T] // 哨兵头节点（虚拟）
	tail *Interval[T] // 哨兵尾节点（虚拟）
	Lock sync.RWMutex
}

// 创建一个新的区间链表，并初始化头尾哨兵节点。
func NewIntervalList[T IntervalValue]() *IntervalList[T] {
	list := &IntervalList[T]{
		head: &Interval[T]{
			StartOffset: -1,
			StopOffset:  -1,
		},
		tail: &Interval[T]{
			StartOffset: math.MaxInt64,
			StopOffset:  math.MaxInt64,
		},
	}
	return list
}

// 获取链表中第一个真实区间（非头哨兵）
func (list *IntervalList[T]) Front() (interval *Interval[T]) {
	return list.head.Next
}

// 在链表尾部追加一个区间（无冲突）
func (list *IntervalList[T]) AppendInterval(interval *Interval[T]) {
	list.Lock.Lock()
	defer list.Lock.Unlock()

	if list.head.Next == nil {
		list.head.Next = interval
	}
	interval.Prev = list.tail.Prev
	if list.tail.Prev != nil {
		list.tail.Prev.Next = interval
	}
	list.tail.Prev = interval
}

// Overlay 用于覆盖插入（直接替换，不考虑版本）指定区间
func (list *IntervalList[T]) Overlay(startOffset, stopOffset, tsNs int64, value T) {
	if startOffset >= stopOffset {
		return
	}
	interval := &Interval[T]{
		StartOffset: startOffset,
		StopOffset:  stopOffset,
		TsNs:        tsNs,
		Value:       value,
	}

	list.Lock.Lock()
	defer list.Lock.Unlock()

	list.overlayInterval(interval)
}

// 直接覆盖插入，不做版本判断（overlay 特性）
func (list *IntervalList[T]) overlayInterval(interval *Interval[T]) {
	p := list.head
	for ; p.Next != nil && p.Next.StopOffset <= interval.StartOffset; p = p.Next {
	}
	q := list.tail
	for ; q.Prev != nil && q.Prev.StartOffset >= interval.StopOffset; q = q.Prev {
	}

	// 处理 interval 左边的剩余（未被覆盖部分）
	if p.Next != nil && p.Next.StartOffset < interval.StartOffset {
		t := &Interval[T]{
			StartOffset: p.Next.StartOffset,
			StopOffset:  interval.StartOffset,
			TsNs:        p.Next.TsNs,
			Value:       p.Next.Value,
		}
		p.Next = t
		if p != list.head {
			t.Prev = p
		}
		t.Next = interval
		interval.Prev = t
	} else {
		p.Next = interval
		if p != list.head {
			interval.Prev = p
		}
	}

	// 处理 interval 右边的剩余（未被覆盖部分）
	if q.Prev != nil && interval.StopOffset < q.Prev.StopOffset {
		t := &Interval[T]{
			StartOffset: interval.StopOffset,
			StopOffset:  q.Prev.StopOffset,
			TsNs:        q.Prev.TsNs,
			Value:       q.Prev.Value,
		}
		q.Prev = t
		if q != list.tail {
			t.Next = q
		}
		interval.Next = t
		t.Prev = interval
	} else {
		q.Prev = interval
		if q != list.tail {
			interval.Next = q
		}
	}
}

// InsertInterval 会插入并和原有区间比较版本号冲突（可分裂、可合并），最终保证无重叠且保留最新。
func (list *IntervalList[T]) InsertInterval(startOffset, stopOffset, tsNs int64, value T) {
	interval := &Interval[T]{
		StartOffset: startOffset,
		StopOffset:  stopOffset,
		TsNs:        tsNs,
		Value:       value,
	}

	list.Lock.Lock()
	defer list.Lock.Unlock()

	value.SetStartStop(startOffset, stopOffset)
	list.insertInterval(interval)
}

// insertInterval 是带冲突处理的插入逻辑核心。
// 会判断是否需要拆分旧区间，保留版本较新的数据区间。
func (list *IntervalList[T]) insertInterval(interval *Interval[T]) {
	prev := list.head
	next := prev.Next

	for interval.StartOffset < interval.StopOffset {
		if next == nil {
			// 如果走到尾部，直接插入
			list.insertBetween(prev, interval, list.tail)
			break
		}

		// 插入区间在 next 之前
		if interval.StopOffset <= next.StartOffset {
			list.insertBetween(prev, interval, next)
			break
		}

		// 插入区间完全在 next 之后
		if next.StopOffset <= interval.StartOffset {
			prev = next
			next = next.Next
			continue
		}

		// 有交集时，根据时间戳判断谁是更新版本
		if interval.TsNs >= next.TsNs {
			// 插入的是更新版本
			if next.StartOffset < interval.StartOffset {
				// 保留左边老数据
				t := &Interval[T]{
					StartOffset: next.StartOffset,
					StopOffset:  interval.StartOffset,
					TsNs:        next.TsNs,
					Value:       next.Value.Clone().(T),
				}
				t.Value.SetStartStop(t.StartOffset, t.StopOffset)
				list.insertBetween(prev, t, interval)
				next.StartOffset = interval.StartOffset
				next.Value.SetStartStop(next.StartOffset, next.StopOffset)
				prev = t
			}
			if interval.StopOffset < next.StopOffset {
				// 保留右边老数据
				next.StartOffset = interval.StopOffset
				next.Value.SetStartStop(next.StartOffset, next.StopOffset)
				list.insertBetween(prev, interval, next)
				break
			} else {
				// 被完全覆盖，跳过 next
				prev.Next = interval
				next = next.Next
			}
		} else {
			// 老数据版本更新，不替换 next
			if interval.StartOffset < next.StartOffset {
				// 插入区间左边部分不重叠，保留为新
				t := &Interval[T]{
					StartOffset: interval.StartOffset,
					StopOffset:  next.StartOffset,
					TsNs:        interval.TsNs,
					Value:       interval.Value.Clone().(T),
				}
				t.Value.SetStartStop(t.StartOffset, t.StopOffset)
				list.insertBetween(prev, t, next)
				interval.StartOffset = next.StartOffset
				interval.Value.SetStartStop(interval.StartOffset, interval.StopOffset)
			}
			if next.StopOffset < interval.StopOffset {
				// 插入区间右边还有剩余，继续处理
				interval.StartOffset = next.StopOffset
				interval.Value.SetStartStop(interval.StartOffset, interval.StopOffset)
			} else {
				// 被覆盖，不插入
				break
			}
		}
	}
}

// 将 interval 插入 a 和 b 中间
func (list *IntervalList[T]) insertBetween(a, interval, b *Interval[T]) {
	a.Next = interval
	b.Prev = interval
	if a != list.head {
		interval.Prev = a
	}
	if b != list.tail {
		interval.Next = b
	}
}

// 获取当前区间链表的长度（不计哨兵）
func (list *IntervalList[T]) Len() int {
	list.Lock.RLock()
	defer list.Lock.RUnlock()

	var count int
	for t := list.head; t != nil; t = t.Next {
		count++
	}
	return count - 1 // 去除头哨兵
}
