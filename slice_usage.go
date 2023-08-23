package main

import "fmt"

// src/runtime/slice.go源码定义：
// type slice struct {
//     array unsafe.Pointer
//     len   int
//     cap   int
// }
// 由上可知slice是一个结构体，保存着指向数组的指针，同时还保存着该切片对应的len和cap
// 注意事项：
//     1. 生成的新切片虽共享底层数组，但是都有着自己的len和cap，即在共享数组上有着自己的访问和修改范围。
//        就如当切片作为参数传递时，如果是非指针传递，那么形参切片会共享底层数组，但是拥有自己的len和cap
//
//     2. 在切片len == cap后进行append操作，会对该切片进行扩容，实质上是生成一个新切片，引用一个扩容后的数组，
//        扩容原则：cap < 1024 --> 2x; cap >= 1024 --> 1.25x
//
//     3.

func main() {
	m1()
}

// m1 方法
func m1() {
	arr := [10]int{}
	// 切片表达式[low:high]表示的是[low, high)区间
	// 切取后的长度为len = high - low
	// 切取后的容积为newCap = oldCap - low
	s := arr[5:6]

	fmt.Printf("len(s) = %d\n", len(s))
	fmt.Printf("cap(s) = %d\n", cap(s))
}

func m2() {

}
