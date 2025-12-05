package main

import (
	"context"
	"fmt"
	"reflect"
	"time"
)

//1. 从 interface{} 变量可以反射出反射对象；

func main() {
	// var x int = 10
	// v := reflect.ValueOf(&x)
	// fmt.Println("v的值是:", v)
	// fmt.Println("v的类型是:", v.Type())
	// fmt.Println("v的kind是:", v.Kind())
	// two()
	// v.Elem().SetInt(20)
	// fmt.Println("v的值是:", v)

	// v := reflect.ValueOf(Add)
	// fmt.Println("v的值是:", v)
	// fmt.Println("v的类型是:", v.Type())
	// fmt.Println("v的kind是:", v.Kind())
	// if v.Kind() == reflect.Func {
	// 	fmt.Println("v是一个函数")
	// 	t := v.Type()
	// 	fmt.Println("t的参数数量是:", t.NumIn())
	// 	fmt.Println("t的返回值数量是:", t.NumOut())
	// 	fmt.Println("t的第一个参数类型是:", t.In(0))
	// 	fmt.Println("t的第一个返回值类型是:", t.Out(0))

	// } else {
	// 	fmt.Println("v不是一个函数")

	// 	return
	// }
	ctx, cancel := context.WithCancel(context.Background())

	go watch(ctx, "监控1")
	go watch(ctx, "监控2")
	go watch(ctx, "监控3")
	time.Sleep(10 * time.Second)
	fmt.Println("10秒后，取消监控")
	cancel()

	time.Sleep(2 * time.Second)
	fmt.Println("2秒后，监控3也被取消")

}

func watch(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println(name, "被取消")
			return
		default:
			fmt.Println(name, "在运行")
			time.Sleep(1 * time.Second)
		}
	}
}

//2. 从反射对象可以获取 interface{} 变量；

func two() {
	v := reflect.ValueOf(1)
	x := v.Interface().(int)
	fmt.Println(x)

}

//3. 要修改反射对象，其值必须可设置；

func Add(a, b int) int {
	return a + b
}
