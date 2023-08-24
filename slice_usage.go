package main

import "fmt"

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
//     3. s := make([]int{}, 10) -> len = 10 & cap = 10
//        s := make([]int{}, 5, 10) -> len = 5 & cap = 10
//     4. arr := [5]int{1, 2, 3, 4, 5}; s := arr[2:4] -> len = 2 & cap = 3
// 底层原理：
//     1. 在切片len == cap后进行append操作，会对该切片进行扩容，实质上是生成一个新切片，引用一个扩容后的新数组，len不变。
//        扩容原则：cap < 1024 --> 2x; cap >= 1024 --> 1.25x
//
//     2. 切片表达式：对数组、字符串或者切片进行部分截取，需要满足0 <= low <= high <= cap
//        i. 简单表达式 -> [low:high]表示的是[low, high)区间
//              生成的切片 len = high - low, cap = oldCap - low
//        ii. 扩展表达式 -> [low:high:max], 其中max >= high
//              为了防止新切片影响共享数组的值，从而影响其他切片，引入max来限定最大容量
//              生成的切片 len = high - low, cap = max - low
//
//     3. 上述新切片共享着底层数组，假设 b := a[low:high]，则最核心的步骤为：
//                       b.array = &a[low]
//        改变新切片的底层数组指针指向的地址，这就是共享数组的原理，b从该新地址开始，b[0～len]取数据。
//
//     4. 当切片作为参数传递时，如果是非指针传递，那么形参切片会共享底层数组，但是拥有自己的len和cap。

func main() {
	m4()
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
	var slice []int
	s1 := append(slice, 1, 2, 3)
	s2 := append(s1, 4)
	fmt.Println(&s1[0] == &s2[0])
}

func m3() {
	orderLen := 5
	order := make([]uint16, 2*orderLen)

	pollorder := order[:orderLen:orderLen]
	lockorder := order[orderLen:][:orderLen:orderLen]

	fmt.Println(len(pollorder))
	fmt.Println(cap(pollorder))
	fmt.Println(len(lockorder))
	fmt.Println(cap(lockorder))
}

func m4() {
	arr := [5]int{1, 2, 3, 4, 5}

	s := arr[2:4]

	for idx, _ := range arr {
		fmt.Println(&arr[idx])
	}
	fmt.Println(&s[0])
}
