package main

import (
	"fmt"
	"math"
	"sync"
	"time"
)

func main() {
	// 指针练习
	fmt.Println("加10")
	var a int = 10
	addTen(&a)
	fmt.Printf("addTen(&a)后 a 的值为: %d\n", a)
	// 2. 实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2
	fmt.Println("乘以2")
	s := []int{1, 2, 3, 4, 5}
	multiplyByTwo(&s)
	fmt.Printf("multiplyByTwo(&s)后 s 的值为: %v\n", s)
	// 3. 编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
	fmt.Println("printOdd测试  start")
	printOdd()
	fmt.Println("printOdd测试  end")

	// 4. 定义一个 Person 结构体，包含 Name 和 Age 字段，再创建一个 Employee 结构体
	// ，组合 Person 结构体并添加 EmployeeID 字段。为
	//  Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。

	// 创建一个 Employee 实例
	e := Employee{
		Person: Person{
			Name: "张三",
			Age:  30,
		},
		EmployeeID: 12345,
	}

	// 调用 PrintInfo 方法
	e.PrintInfo()

	// chainTest测试
	fmt.Println("chainTest测试  start")
	chainTest()
	fmt.Println("chainTest测试  end")

	// counterTest测试
	fmt.Println("counterTest测试  start")
	counterTest()
	fmt.Println("counterTest测试  end")

}

// 指针练习
// 1题目 ：编写一个Go程序，
// 定义一个函数，该函数接收一个整数指针作为参数，
// 在函数内部将该指针指向的值增加10，然后在主函数中调用该函数并输出修改后的值。
func addTen(p *int) {
	*p += 10
}

//2.实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2

func multiplyByTwo(s *[]int) {

	for i := range *s {
		(*s)[i] *= 2
	}
}

//Goroutine
//编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。

func printOdd() {

	go func() {
		for i := 1; i <= 10; i += 2 {
			fmt.Printf("奇数携程: %d", i)
		}
	}()
	go func() {
		for i := 2; i <= 10; i += 2 {
			fmt.Printf("偶数携程: %d", i)
		}
	}()
	//等待所有协程完成（简单起见，使用睡眠）
	//在实际项目中，应该使用sync.WaitGroup
	time.Sleep(time.Second)
}

// 设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
type Task func()

func TaskScheduler(tasks []Task) {
	var wg sync.WaitGroup
	for _, task := range tasks {
		wg.Add(1)
		go func(t Task) {
			defer wg.Done()
			start := time.Now()
			t()
			elapsed := time.Since(start)
			fmt.Printf("任务执行时间: %v\n", elapsed)
		}(task)
	}
	wg.Wait()
}

//定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法、
// 然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。在主函数中，
// 创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

//：使用组合的方式创建一个 Person 结构体，
// 包含 Name 和 Age 字段，再创建一个 Employee 结构体
// ，组合 Person 结构体并添加 EmployeeID 字段。为
//  Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。

type Person struct {
	Name string
	Age  int
}

type Employee struct {
	Person
	EmployeeID int
}

func (e Employee) PrintInfo() {
	fmt.Printf("姓名: %s, 年龄: %d, 员工ID: %d\n", e.Name, e.Age, e.EmployeeID)
}

//题目 ：编写一个程序，使用通道实现两个协程之间的通信。
// 一个协程生成从1到10的整数，并将这些整数发送到通道中，
// 另一个协程从通道中接收这些整数并打印出来。

func chainTest() {

	//1.创建一个通道
	ch := make(chan int, 10)
	//一个协程生成从1到10的整数，并将这些整数发送到通道中，
	go func() {
		for i := 1; i <= 10; i++ {
			ch <- i
		}
		//发送完成后关闭通道
		close(ch)

	}()
	//另一个协程从通道中接收这些整数并打印出来。
	go func() {
		//使用range循环接收，当通道关闭且数据接收完毕时自动退出
		for v := range ch {
			fmt.Println(v)
		}
	}()

	//等待所有协程完成（简单起见，使用睡眠）
	//在实际项目中，应该使用sync.WaitGroup
	time.Sleep(time.Second)
}

var mu sync.Mutex

// 编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
func counterTest() int {
	var sum int

	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 1000; j++ {
				mu.Lock()
				sum++
				mu.Unlock()
			}
		}()

	}
	//等待所有协程完成（简单起见，使用睡眠）
	//在实际项目中，应该使用sync.WaitGroup
	time.Sleep(time.Second)
	return sum
}
