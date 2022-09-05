// Copyright © 2022 UCloud. All rights reserved.

package main

import (
	"fmt"
	"testing"
)

var ch chan int

func test() int {
	ch = make(chan int)
	go func() {
		// for {
		fmt.Println(<-ch) //
		fmt.Println("hello")
		// }
		fmt.Println("aaaa")
	}()
	// 不阻塞，那go func()不会异常退出吗？
	// 协程并不是函数，不会因为这个函数的退出而退出
	// test()启动一个deadloop子协程,这个会在主协程main结束后被强制退出
	return 0
}
func TestGoRoutine(t *testing.T) {
	c := test()
	fmt.Println("c", c)
	ch <- 10 // 塞一个数
}
