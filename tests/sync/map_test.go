package sync

import (
	"fmt"
	"sync"
	"testing"
)

func TestBasic(t *testing.T) {
	//初始化
	m := sync.Map{}
	//加载
	m.Store("name", "duyuqing")
	value, ok := m.Load("name")
	if ok {
		fmt.Println(value)
	}
	value1, loaded := m.LoadAndDelete("name")
	if loaded {
		fmt.Println(value1)
	}
	actual, loaded := m.LoadOrStore("name", "duyuqing")
	if !loaded {
		fmt.Println(actual)
	}
	m.Delete("name")

	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("name-%v", i)
		m.Store(key, i)
	}
	m.Range(func(key, value any) bool {
		fmt.Println(key, ":", value)
		return true
	})
}

func Print1() {
	fmt.Println("aaa")
}

func Print2() {
	fmt.Println("bbb")
}

func TestAdvanced(t *testing.T) {
	m := sync.Map{}
	m.Store("f1", Print1)
	m.Store("f2", Print2)
	m.Range(func(key, value any) bool {
		f, ok := value.(func())
		if ok {
			f()
		}
		return true
	})
}

func TestComparable(t *testing.T) {
	m := sync.Map{}
	m.Store("f1", []int{1, 2, 3})
	m.Store("f2", []int{4, 5, 6})
	m.Range(func(key, value any) bool {
		fmt.Println(key, ":", value)
		return true
	})
}

type Array[T any] struct {
	data []T
}

func (a *Array[T]) Display() {
	for i := 0; i < len(a.data); i++ {
		fmt.Println(a.data[i])
	}
}

func TestGenericity(t *testing.T) {
	a := Array[int]{data: []int{1, 2, 3}}
	a.data = append(a.data, 6)
	m := sync.Map{}
	m.Store("array", a)
	m.Range(func(key, value any) bool {
		fmt.Println(key, ":", value)
		return true
	})
}
