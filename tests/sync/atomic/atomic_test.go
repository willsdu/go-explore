package atomic

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var count = atomic.Int32{}

func TestAtomicCompetion(t *testing.T) {
	start := time.Now()
	wg := &sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			count.Add(1)
			fmt.Printf("count:%v\n", count.Load())
		}(wg)
	}
	wg.Wait()
	fmt.Printf("Last:%v;Milliseconds:%v", count.Load(), time.Since(start).Microseconds())
}

var number = 0

func TestCompetion(t *testing.T) {
	start := time.Now()
	wg := &sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			number += 1
			fmt.Printf("count:%v\n", number)
		}(wg)
	}
	wg.Wait()
	fmt.Printf("Last:%v;Milliseconds:%v", number, time.Since(start).Microseconds())
}
