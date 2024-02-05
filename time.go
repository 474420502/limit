package limit

import (
	"sync"
	"time"
)

// TimeLimit结构体包含需要限频的映射表,以及时间间隔
// 使用sync.Mutex保证对字典的线程安全访问
type TimeLimit[T comparable] struct {
	mu   sync.Mutex
	dict map[T]struct{}
	dur  time.Duration
}

// NewTimeLimit构造函数,接收限频的时间间隔
// 并初始化内部字典和间隔字段
func NewTimeLimit[T comparable](dur time.Duration) *TimeLimit[T] {
	return &TimeLimit[T]{
		dict: make(map[T]struct{}),
		dur:  dur,
	}
}

// WithTime方法用于更新限频的时间间隔
func (tup *TimeLimit[T]) WithTime(dur time.Duration) *TimeLimit[T] {
	tup.mu.Lock()
	defer tup.mu.Unlock()
	tup.dur = dur
	return tup
}

// ExceedsLimit方法检查传入值是否是一个新的值
// 首先会查询字典,如果存在则表示在间隔内已经访问过
// 否则将其添加到字典,并启动一个定时器在间隔后删除
// ExceedsLimitDo 函数如果返回true. 只有超过定时器时间才能再次触发
func (tup *TimeLimit[T]) ExceedsLimit(v T, ExceedsLimitDo func() bool) {
	tup.mu.Lock()
	defer tup.mu.Unlock()

	if _, ok := tup.dict[v]; ok {
		return
	}

	if ExceedsLimitDo() {
		tup.dict[v] = struct{}{}
		time.AfterFunc(tup.dur, func() {
			tup.mu.Lock()
			defer tup.mu.Unlock()
			delete(tup.dict, v)
		})
	}

	return
}
