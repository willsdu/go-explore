# Golang sync.noCopy类型

> 转载自：[golang sync.noCopy 类型 —— 初探 copylocks 与 empty struct](https://www.cnblogs.com/leozmm/p/go_nocopy_copylocks_empty_struct.html)

# 问题引入

学习golang（v1.16）的 WaitGroup 代码时，看到了一处奇怪的用法，见下方类型定义：

```golang
    type WaitGroup struct {
        noCopy noCopy
        ...
    }
```

这里，有个奇怪的“`noCopy`”类型，顾名思义，这个应该是某种“不可复制”的意思。下边是noCopy类型的定义：

```golang
    // noCopy may be embedded into structs which must not be copied
    // after the first use.
    //
    // See https://golang.org/issues/8005#issuecomment-190753527
    // for details.
    // 对应github链接：https://github.com/golang/go/issues/8005#issuecomment-190753527
    type noCopy struct {}
    // Lock is a no-op used by -copylocks checker from `go vet`
    func (*noCopy) Lock{}
    func (*noCopy) Unlock{}
    // 以上 Lock 和 Unlock 方法属于 Locker 接口类型的方法集，见 sync/mutex.go
```

这里有2点比较特别：

1. noCopy 类型是空 struct
2. noCopy 类型实现了两个方法： Lock 和 Unlock，而且都是空方法（no-op）。注释中有说，这俩方法是给 go vet 的 copylocks 检测器用的

也就是说，这个 noCopy 类型和它的方法集，没有任何实质的功能属性。那么它是用来做什么的呢？

# 动手试试

从类型定义，以及实现的lock方法的注释可以看出，noCopy 是为了实现对不可复制类型的限制。这个限制如何起作用呢？参考注释中给出的 issuecomment 链接，在[Russ Cox](https://github.com/rsc) 的评论中，看的这么一句：

> A package can define:
>
> ```golang
> type noCopy struct{}
> func (*noCopy) Lock() {}
> ```
>
> and then put a `noCopy noCopy` into any struct that must be flagged by vet.

原来这个`noCopy`的用处，是为了让被嵌入的container类型，在用`go vet`工具进行copylock check时，能被检测到。

我写了一段代码试了下：

```golang
    // file: main.go
    package main
    import "fmt"
    type noCopy struct{}
    func (*noCopy) Lock()   {}
    func (*noCopy) Unlock() {}
    type cool struct {
    	Val int32
    	noCopy
    }

    func main() {
    	c1 := cool{Val:10,}
    	c2 := c1                // <- 赋值拷贝
    	c2.Val = 20
    	fmt.Println(c1, c2)     // <- 传参拷贝
    }
```

然后，我先用vet工具检查了一下：

```shell
    leo@leo-MBP % go vet main.go
    # command-line-arguments
    ./main.go:14:8: assignment copies lock value to c2: command-line-arguments.cool
    ./main.go:16:14: call of fmt.Println copies lock value: command-line-arguments.cool
    ./main.go:16:18: call of fmt.Println copies lock value: command-line-arguments.cool
```

上边的输出可以看到，在代码标记出来的两处位置，vet打印了“copy lock value”的提示。

# 查找资料

试着查了一下这个提示的相关信息，发现这一篇博文：[Detect locks passed by value in Go](https://medium.com/golangspec/detect-locks-passed-by-value-in-go-efb4ac9a3f2b)

同时，用`go tool vet help copylocks`命令可以查看 vet 对 copylocs 分析器的介绍：

> copylocks: check for locks erroneously passed by value
>
> Inadvertently copying a value containing a lock, such as sync.Mutex or
> sync.WaitGroup, may cause both copies to malfunction. Generally such
> values should be referred to through a pointer.

原来，vet 工具的 copylocks 检测器有这么一个功能：检测带锁类型（如 `sync.Mutex`） 的错误复制使用，这种不当的复制，会引发死锁。

其实不仅仅是`sync.Mutex`类型会这样，所有需要用到Lock和Unlock方法的类型，即 lock type，都有这种 **“错误复制引发死锁”** 的隐患。

所以，我们在上边测试的代码中定义的`noCopy`类型，实现了`Lock`和`Unlock`方法，使得 noCopy 成了一个 lock type，目的就是为了能利用 vet 的 copylocks 分析器对 copy value 的检测能力。

> 岔个题
> 虽然上边的测试代码，在用 go vet 检测时给出了提示信息，但是这并不是警告，相应代码没有语法错误，仍然是可执行的，run 一下试试：
>
> ```shell
> leo@leo-MBP % go run main.go
> {10 {}} {20 {}}
> ```
>
> 嵌入了 noCopy 类型的 cool 类型，在被强行复制之后，依然可以运行。noCopy 这种设计的意义，在于防范不当的 copylocks 发生，且这种防范不是强制的，依靠开发者自行检测。

# 空struct

好，明白了 `noCopy` 的存在的意义，接下来探究一下 `noCopy` 为什么要设计成空 struct 类型。

先上结论：使用空 struct 是出于性能考虑。

```golang
    package main

    import (
    	"fmt"
    	"unsafe"
    )

    type cool struct{}

    func main() {
    	c := cool{}
    	fmt.Println(unsafe.Sizeof(c)) // -> print 0
    }
```

如上所示，空 struct 类型的值不占用内存空间，所以在性能上更有优势。

# 总结

综合来看，noCopy 空 struct 类型，结合了 vet 工具对 copylocks 检测的支持，以及空 struct 对性能的优化，用在 **“标记不可复制类型”** 的场景下，是比较巧妙的设计。

# 参考

[Detect locks passed by value in Go](https://medium.com/golangspec/detect-locks-passed-by-value-in-go-efb4ac9a3f2b)
[The empty struct](https://dave.cheney.net/2014/03/25/the-empty-struct#comment-2815)
[Go 空结构体 struct{} 的使用](https://geektutu.com/post/hpg-empty-struct.html)