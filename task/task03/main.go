package main

import (
	"fmt"
)

func main() {
	fmt.Println("hello world")

	v := [3]int{1, 2, 3}
	s := v[0:2]
	s[0] = 1000
	v[1] = 2000

	fmt.Println("s: %v", s)
	fmt.Println("v: %v", v)

	updateByIdex(&v, 1, 100)
	fmt.Println(v)
}

func updateByIdex(v *[3]int, index int, value int) {
	if index >= 0 && index < len(v) {
		v[index] = value
	}
}
