package limit

import (
	"log"
	"testing"
	"time"
)

func TestCaseFreq(t *testing.T) {
	freq := NewFrequencyLimit()
	for i := 0; i < 1000; i++ {
		freq.Put(10)
		time.Sleep(time.Second)
		log.Println(freq.GetFrequency())
	}

}
