package main

import (
	"fmt"
)

// src/runtime/slice.go源码定义：
// type slice struct {
//     array unsafe.Pointer
//     len   int
//     cap   int
// }
// 由上可知slice是一个结构体，保存着指向数组的指针，同时还保存着该切片对应的len和cap
// 声明及初始化：
//     1. var s []int -> nil
//     2. s := []int{} -> len = 0 & cap = 0
//        s := []int{1, 2, 3} -> len = 3 & cap = 3
//     3. s := make([]int, 10) -> len = 10 & cap = 10
//        s := make([]int, 5, 10) -> len = 5 & cap = 10
//        s := make([]int, []int{1, 2, 3}...) -> len
//     4. arr := [5]int{1, 2, 3, 4, 5}; s := arr[2:4] -> len = 2 & cap = 3
// 底层原理：
//     1. 在切片len == cap后进行append操作，会对该切片进行扩容，实质上是生成一个新切片，引用一个扩容后的新数组，len不变。
//        扩容原则（可参照slice.go#growslice方法）: 设新增元素后的总容量为aftercap
//            aftercap > 2 * cap --> aftercap
//            cap < 256 --> 2x;
//            cap >= 256 --> 1.25x:
//                        for newcap < aftercap {
//                            newcap = 0.25 * newcap + 0.75 * threshold
//                        }
//                   go1.18.3中threshold为256；
//                   只扩容一次的情况下，当slice容量特别大时，加上的常量可以忽略不计，即总体上扩容1.25x
//
//     2. 切片表达式：对数组、字符串或者切片进行部分截取，需要满足0 <= low <= high <= cap
//        i. 简单表达式 -> [low:high]表示的是[low, high)区间
//              生成的切片 len = high - low, cap = oldCap - low
//        ii. 扩展表达式 -> [low:high:max], 其中max >= high
//              为了防止新切片影响共享数组的值，从而影响其他切片，引入max来限定最大容量
//              生成的切片 len = high - low, cap = max - low
//        简单表达式中low和high都可省略，拓展表达式中由于max和high强绑定，故只有low可省略
//
//     3. 上述新切片共享着底层数组，假设 b := a[low:high]，则最核心的步骤为：
//                       b.array = &a[low]
//        改变新切片的底层数组指针指向的地址，这就是共享数组的原理，b从该新地址开始，b[0～len]取数据。
//
//     4. 当切片作为参数传递时，如果是非指针传递，那么形参切片会共享底层数组，但是拥有自己的len和cap。

// m1 方法
func m1() {
	s := []int{1}
	tmp1 := append(s, 2)
	tmp2 := append(tmp1, 3)

	fmt.Printf("s\tadd -> %p, len -> %d, cap -> %d\n", &s, len(s), cap(s))
	fmt.Printf("tmp1\tadd -> %p, len -> %d, cap -> %d\n", &tmp1, len(tmp1), cap(tmp1))
	fmt.Printf("tmp2\tadd -> %p, len -> %d, cap -> %d\n", &tmp2, len(tmp2), cap(tmp2))
}

func m2() {
	s := make([]int, 1, 3)
	s1 := append(s, 1)
	s2 := append(s1, 2)
	// 还未扩容，s1和s2共享底层数组，但是有自己的len和cap
	fmt.Printf(
		"s1 -> %+v, len(s1) -> %d, cap(s1) -> %d, \n"+
			"s2 -> %+v, len(s2) -> %d, cap(s2) -> %d, \n"+
			"share -> %v \n",
		s1, len(s1), cap(s1),
		s2, len(s2), cap(s2),
		&s1[0] == &s2[0])

	fmt.Println("-------------------------------------------")
	s3 := append(s2, 3)
	// 由于扩容，底层的数组都不是同一个了
	fmt.Printf(
		"s2 -> %+v, len(s2) -> %d, cap(s2) -> %d, \n"+
			"s3 -> %+v, len(s3) -> %d, cap(s3) -> %d, \n"+
			"share -> %v \n",
		s2, len(s2), cap(s2),
		s3, len(s3), cap(s3),
		&s2[0] == &s3[0])
}

func m3() {
	orderLen := 5
	order := make([]uint16, 2*orderLen)

	pollorder := order[:orderLen:orderLen]
	lockorder := order[orderLen:][:orderLen:orderLen]

	fmt.Printf("pollorder len -> %d, cap -> %d\n", len(pollorder), cap(pollorder))
	fmt.Printf("lockorder len -> %d, cap -> %d\n", len(lockorder), cap(lockorder))
}

// m4 修改指针的指向
func m4() {
	arr := [5]int{1, 2, 3, 4, 5}
	s := arr[2:4]

	for idx, _ := range arr {
		fmt.Println(&arr[idx])
	}
	// 发现切片的指针指向了arr[2]的地址
	fmt.Println(&s[0])
}
