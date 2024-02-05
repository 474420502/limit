package limit

import (
	"log"
	"testing"
	"time"
)

func TestCaseFreq(t *testing.T) {
	freq := NewFrequencyLimit[Float64]()

	for i := 0; i < 100; i++ {
		time.Sleep(time.Second)
		freq.Put(Float64(1))
		log.Println(freq.GetFrequency())
	}

}
