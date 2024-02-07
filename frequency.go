package limit

import (
	"sync"
	"time"
)

type countingUnit struct {
	Counted float64
	Time    time.Time
}

type FrequencyLimit struct {
	total        float64
	countedUints List[*countingUnit]

	countedMaxSize int

	lock sync.Mutex
}

func NewFrequencyLimit() *FrequencyLimit {
	return &FrequencyLimit{
		total:          0,
		countedMaxSize: 64,
	}
}

func (freq *FrequencyLimit) With(countedMaxSize int) {
	freq.countedMaxSize = countedMaxSize
	for freq.countedUints.size > freq.countedMaxSize {
		freq.total -= freq.countedUints.RemoveBack().Counted
	}
}

func (freq *FrequencyLimit) Put(counted float64) {
	freq.lock.Lock()
	defer freq.lock.Unlock()

	if counted == 0 {
		return
	}

	cunit := &countingUnit{
		Counted: counted,
	}
	defer func() { cunit.Time = time.Now() }()

	freq.countedUints.PutHead(cunit)
	freq.total += cunit.Counted

	if freq.countedUints.size > freq.countedMaxSize {
		freq.total -= freq.countedUints.RemoveBack().Counted
	}

}

func (freq *FrequencyLimit) getDuration() time.Duration {
	dur := time.Since(freq.countedUints.tail.value.Time)
	return dur
}

// 默认是秒
func (freq *FrequencyLimit) GetFrequency() float64 {
	freq.lock.Lock()
	defer freq.lock.Unlock()

	if freq.countedUints.size == 0 {
		return 0.0
	}

	return freq.total / freq.getDuration().Seconds()
}

// 非默认秒 频率计算, do 输入的是时间差
func (freq *FrequencyLimit) GetFrequencyWith(do func(freqTime time.Duration) float64) float64 {
	freq.lock.Lock()
	defer freq.lock.Unlock()

	if freq.countedUints.size == 0 {
		return 0.0
	}

	return freq.total / do(freq.getDuration())
}
