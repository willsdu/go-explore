

https://blog.csdn.net/u013911096/article/details/139195294

https://blog.csdn.net/qq_29647593/article/details/136166769

https://github.com/golang/go/issues/56603

https://www.jianshu.com/p/f3230d5f69df



# 解读: `_ [0]*T`

```go
// A Pointer is an atomic pointer of type *T. The zero value is a nil *T.
type Pointer[T any] struct {
	// Mention *T in a field to disallow conversion between Pointer types.
	// See go.dev/issue/56603 for more details.
	// Use *T, not T, to avoid spurious recursive type definition errors.
	_ [0]*T

	_ noCopy
	v unsafe.Pointer
}
```

看这个issue的代码

```go
package main

import (
	"math"
	"sync/atomic"
)

type small struct {
	small [64]byte
}

type big struct {
	big [math.MaxUint16 * 10]byte
}

func main() {
	a := atomic.Pointer[small]{}
	a.Store(&small{})

	b := atomic.Pointer[big](a) // type conversion
	big := b.Load()

	for i := range big.big {
		big.big[i] = 1
	}
}
```

如果不加这个, 会访问不该访问的内存，导致奔溃