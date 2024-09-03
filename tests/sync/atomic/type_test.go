package atomic

import (
	"sync/atomic"
	"testing"
)

func TestBool(t *testing.T) {
	a := atomic.Bool{}
	a.Store(true)
	t.Log(a.Load())
	t.Log(a.Swap(false), a.Load())
	t.Log(a.CompareAndSwap(true, false))
	t.Log(a.CompareAndSwap(false, true))
}
