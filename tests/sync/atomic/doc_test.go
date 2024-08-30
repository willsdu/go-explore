package atomic

import (
	"fmt"
	"sync/atomic"
	"testing"
	"unsafe"
)

var x int
var y = 10
var z = 20

func TestCompare(t *testing.T) {
	px := unsafe.Pointer(&x)
	py := unsafe.Pointer(&y)
	pz := unsafe.Pointer(&z)
	fmt.Printf("%p:%p\n", px, py)
	r := atomic.CompareAndSwapPointer(&px, pz, pz)
	t.Log(r)
	fmt.Printf("%p:%p\n", px, py)
}

type Ptr struct {
	x unsafe.Pointer
}

func TestLoadPointer(t *testing.T) {
	p := Ptr{
		x: unsafe.Pointer(&y),
	}
	a := atomic.LoadPointer(&p.x)
	s := (*int)(a)
	t.Log(*s)
}

func TestLoadPointer2(t *testing.T) {
	p := unsafe.Pointer(&x)
	a := atomic.LoadPointer(&p)
	t.Log(a)
	t.Log(unsafe.Sizeof(x))
	fmt.Printf("%16p\n", &p)
}
