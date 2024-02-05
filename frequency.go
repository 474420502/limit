package limit

import (
	"sync"
	"time"
)

type IFloat64 interface {
	CountedFloat() float64
}

type countingUnit struct {
	Counted float64
	Time    time.Time
}

type FrequencyLimit struct {
	total        float64
	countedUints List[*countingUnit]

	minDuration time.Duration

	lock sync.Mutex
}

func NewFrequencyLimit() *FrequencyLimit {
	return &FrequencyLimit{
		total:       0,
		minDuration: time.Second * 4,
	}
}

func (freq *FrequencyLimit) Put(counted float64) {
	freq.lock.Lock()
	defer freq.lock.Unlock()

	cunit := &countingUnit{
		Counted: counted,
		Time:    time.Now(),
	}
	freq.countedUints.PutHead(cunit)
	freq.total += cunit.Counted

	if freq.countedUints.head == freq.countedUints.tail {
		return
	}

	tail := freq.countedUints.tail
	head := freq.countedUints.head

	for tail != head {
		if head.value.Time.Sub(tail.value.Time) >= freq.minDuration {
			next := freq.countedUints.TruncateNodeNext(tail)
			for _, v := range next.ToTailValues() {
				freq.total -= v.Counted
			}
			return
		}

		tail = tail.prev
	}

}

func (freq *FrequencyLimit) getDuration() time.Duration {
	return freq.countedUints.head.value.Time.Sub(freq.countedUints.tail.value.Time)
}

// 默认是秒
func (freq *FrequencyLimit) GetFrequency() float64 {
	freq.lock.Lock()
	defer freq.lock.Unlock()

	if freq.countedUints.head == freq.countedUints.tail {
		return 0.0
	}

	return freq.total / freq.getDuration().Seconds()
}

// 非默认秒 频率计算, do 输入的是时间差
func (freq *FrequencyLimit) GetFrequencyWith(do func(freqTime time.Duration) float64) float64 {
	freq.lock.Lock()
	defer freq.lock.Unlock()

	if freq.countedUints.head == freq.countedUints.tail {
		return 0.0
	}

	return freq.total / do(freq.getDuration())
}
