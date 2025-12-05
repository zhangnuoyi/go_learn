package main

import (
	"fmt"
	"reflect"

	"go_learn/05go/advanced/concurrency"
)

func main() {

	name := "moobb"

	fmt.Println("name 的类型是:", reflect.TypeOf(name))
	fmt.Println("name 的值是:", reflect.ValueOf(name))
	// 调用 CreatedAtGo 函数
	//主死协从
	concurrency.CreatedAtGo()

}
