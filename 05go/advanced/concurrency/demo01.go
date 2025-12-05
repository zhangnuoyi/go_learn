package concurrency

import (
	"fmt"
	"sync"
)

// Java python php 基于多线程，主要存在内存占用高，线程之间的切换成本高
// Golang 基于协程 内存占有少（2K），协程之间的切换成本低，go语言只有协程
var wg sync.WaitGroup

func CreatedAtGo() {
	// // 创建一个新的协程
	// go func() {
	// 	fmt.Println("Hello, World!")
	// }()

	for i := range 100 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fmt.Println(i)
		}(i)
	}
	// 主协程继续执行
	fmt.Println("Main goroutine is running")
	wg.Wait()
}
