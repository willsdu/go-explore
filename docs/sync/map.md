# sync.map的使用及原理

[toc]

## 一、介绍

Go的原生map是非线程安全的，如果想要并发读写map，需要加锁处理。但是加锁处理会影响到性能问题。于是，Go官方给出了一个并发安全的map: sync.map，特点如下：

1. 并发安全：读写分离，同步时加锁，保证并发访问安全
1. 内存安全：主动释放指针内存
1. 性能优：使用读缓存map，减少了加锁次数。

## 二、使用

先看看怎么使用再说

### 2.1 读

> func (m *Map) Load(key any) (value any, ok bool)

```
m:=sync.map{}
name, ok:=m.Store("name")   
if ok {
	fmt.Println(name)
}
```

### 2.2 写

> func (m *Map) Store(key, value any)

```ok
m:=sync.map{}
m.Store("go","ok")
```

### 2.3 删

> func (m *Map) Delete(key any)

```
m:=sync.map{}
m.Store("name")
```



### 2.4 其他

#### 2.4.1 LoadOrStore

```go
func (m *Map) LoadOrStore(key, value any) (actual any, loaded bool)
```

LoadOrStore逻辑如下：

- key存在：actual=key对应的值，loaded=true

- key不存在：actual=value，loaded=false

#### 2.4.2 LoadAndDelete

```go
func (m *Map) LoadAndDelete(key any) (value any, loaded bool)
```

LoadAndDelete的逻辑如下：

- key存在：value=key对应的值，loaded=true

- key不存在：value=nil，loaded=false

## 三、原理





## 四、源码







# 参考

https://zhuanlan.zhihu.com/p/620632865

https://blog.51cto.com/u_16099271/6806907

https://blog.csdn.net/qq_37102984/article/details/128154924