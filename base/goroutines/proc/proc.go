// Copyright © 2022 UCloud. All rights reserved.

package main

import (
	"fmt"
	"sync"
)

func demo() int {
	var aa = 10
	defer func() {
		fmt.Println("sub func", aa)
		aa = 100
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func(i int) {
		defer func() {
			aa = 10000
			fmt.Println("goroutine", aa)
			wg.Done()
		}()
		aa = 100
	}(aa)
	wg.Wait()
	return aa
}
func main() {
	fmt.Println(demo())
	return
	// for {
	// 	var x *int32
	// 	fmt.Println("x", x)
	// 	*x = 0
	// }
}
