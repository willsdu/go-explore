package unsafe

import (
	"testing"
	"unsafe"
)

func TestPointer(t *testing.T) {
	x := unsafe.Pointer(new(any))
	t.Log(x)
	t.Log(x == nil)
}
