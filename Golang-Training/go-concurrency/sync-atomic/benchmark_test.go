package syncatomic_test

import (
	"sync"
	"sync/atomic"
	"testing"
)

const iterations = 10000000

func BenchmarkMutex(t *testing.B) {
	var wg sync.WaitGroup
	var mutex sync.Mutex
	sum := 0

	wg.Add(iterations)

	for i := 0; i < iterations; i++ {
		go func(val int) {
			defer wg.Done()
			mutex.Lock()
			sum += val
			mutex.Unlock()
		}(i)
	}

	wg.Wait()
}

func BenchmarkAtomic(t *testing.B) {
	var wg sync.WaitGroup
	sum := int64(0)

	wg.Add(iterations)

	for i := 0; i < iterations; i++ {
		go func(val int) {
			defer wg.Done()
			atomic.AddInt64(&sum, int64(val))
		}(i)
	}

	wg.Wait()
}
