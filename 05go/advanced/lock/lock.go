package main

import (
	"fmt"
	"sync"
)

var tatal int32 = 0
var wg sync.WaitGroup
var lock sync.Mutex

func add() {
	for i := range 100 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			// atomic.AddInt32(&tatal, 1)
			lock.Lock()
			tatal++
			lock.Unlock()
		}(i)
	}
}

func sub() {
	for i := range 100 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			// atomic.AddInt32(&tatal, 1)
			lock.Lock()
			tatal--
			lock.Unlock()
		}(i)
	}
}

func main() {
	add()
	sub()
	wg.Wait()
	fmt.Println(tatal)
}
