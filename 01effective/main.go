package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"runtime"
	"time"
)

func main() {
	// fmt.Println("Hello, World!")
	// //if 测试
	// fmt.Println("---------if start---------")
	// fmt.Println(ifTest(10))
	// fmt.Println(ifTest(-10))
	// fmt.Println("---------if end---------")
	// //for 测试
	// fmt.Println("---------for start---------")
	// forTest1(10)
	// forTest2(10)
	// forTest3(10)
	// fmt.Println("---------for end---------")
	// //switch 测试
	// switchTest()
	// switchTest2(2)
	// switchTest3(10)
	// switchTest3("hello")
	// switchTest3(10.0)
	// //function 测试
	// fmt.Println("---------function start---------")
	// a, b := swap(10, "hello")
	// fmt.Printf("%v %v\n", a, b)
	// fmt.Println("---------function end---------")

	// // 文件读取测试
	// fmt.Println("---------readFile start---------")
	// content, err := readFile("../go.mod")
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(content)
	// }
	// fmt.Println("---------readFile end---------")

	// // 命名结果参数测试
	// fmt.Println("---------sumAndB start---------")
	// fmt.Println(sumAndB(10, 20))
	// fmt.Println("---------sumAndB end---------")

	// // defer 测试
	// deferTest()
	// // newTest
	// newTest()

	// typeAssertion()

	// fmt.Println(y)

	selectTest()
	divideError(10, 2)
	result, err := divideError(10, 0)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}
}

func ifTest(a int) int {
	if a > 0 {
		return a
	}
	return -a
}

func forTest1(count int) {
	for i := 0; i < count; i++ {
		fmt.Println(i)
	}
}

func forTest2(count int) {
	i := 0
	for i < count {
		fmt.Println(i)
		i++
	}
}

func forTest3(count int) {
	i := 0
	for {
		fmt.Println(i)
		i++
		if i >= count {
			break
		}
	}

}

// switch 表达式可以是任何类型
func switchTest() {
	fmt.Println("---------switch start---------")
	switch os := runtime.GOOS; os {
	case "linux":
		fmt.Println("linux")
	case "windows":
		fmt.Println("windows")
	default:
		fmt.Println("default")
	}
	fmt.Println("---------switch end---------")
}

func switchTest2(i int) {
	fmt.Println("---------switch2 fallthrough start---------")
	switch {
	case i < 5:
		fmt.Println("i < 5")
		fallthrough
	case i < 10:
		fmt.Println("i < 10")
	default:
		fmt.Println("default")
	}
	fmt.Println("----------switch2 fallthrough  end---------")
}
func switchTest3(i interface{}) {
	fmt.Println("---------switch3 type  start---------")
	switch i.(type) {
	case int:
		fmt.Printf("int  %v\n", i)
	case string:
		fmt.Printf("string  %v\n", i)
	default:
		fmt.Printf("default  %v\n", i)
	}
	fmt.Println("----------switch3 type  end---------")
}

// 2个对象交换
func swap(a, b interface{}) (interface{}, interface{}) {
	fmt.Printf("swap befor %v %v\n", a, b)
	defer fmt.Printf("swap after  %v %v\n", b, a)
	return b, a
}

// 文件读取，返回文件内容和错误信息
func readFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "没有读取到文件内容", err
	}
	return string(data), nil
}

// 命名结果参数
func sumAndB(a, b int) (sum int) {
	sum = a + b
	return sum
}

// defer 测试
func deferTest() {
	fmt.Println("---------defer start---------")
	defer fmt.Println("defer 1")
	defer fmt.Println("defer 2")
	fmt.Println("defer 3")
	fmt.Println("---------defer end---------")
}

// newTest
func newTest() {
	fmt.Println("---------newTest start---------")
	s := new(string)
	fmt.Printf("%v %T\n", s, *s)
	*s = "hello"
	fmt.Printf("%v %T\n", *s, *s)
	i := new(int)
	fmt.Printf("%v %T\n", i, i)
	*i = 10
	fmt.Printf("%v %T\n", *i, *i)

	fmt.Println("---------newTest end---------")
}

//Go 没有内置的构造函数，但你可以通过函数来实现构造函数的功能

type Point struct {
	X int
	Y int
}

func NewPoint(x, y int) *Point {
	return &Point{x, y}
}

var y = initX()

func initX() int {
	return 10
}

// init
func init() {
	fmt.Println("init1")
	y = 20
}

func init() {
	fmt.Println("init2")
	y = 30
}

func init() {
	fmt.Println("init3")
	y = 40
}

// 方法的接收器可以是值类型或指针类型。
// 计算两点之间的距离
func (p Point) Distance(q Point) float64 {
	dx := p.X - q.X
	dy := p.Y - q.Y
	return math.Sqrt(float64(dx*dx + dy*dy))
}

// 指针接收
// 缩放点的坐标
func (p *Point) Scale(factor int) {
	p.X *= factor
	p.Y *= factor
}

type Reader interface {
	Read(p []byte) (n int, err error)
}

type Writer interface {
	Write(p []byte) (n int, err error)
}
type File struct {
	// 文件相关字段
}

// File 实现了 Reader 接口
func (f *File) Read(p []byte) (n int, err error) {
	// 实现读取逻辑
	return 0, nil
}

// File 实现了 Writer 接口
func (f *File) Write(p []byte) (n int, err error) {
	// 实现写入逻辑
	return 0, nil
}

//类型断言

var i interface{} = 10

func typeAssertion() {

	if v, ok := i.(string); ok {
		fmt.Printf("i is string %v\n", v)
	} else {
		fmt.Printf("i is not int %v\n", i)
	}
}

func channels() {
	fmt.Println("---------channels start---------")
	ch := make(chan int)
	go func() {
		ch <- 10
	}()
	v := <-ch
	fmt.Println("获取通道的值 v = ", v)
	fmt.Println("---------channels end---------")
}

func channels2() {
	fmt.Println("---------channels start---------")
	ch := make(chan int, 10)
	go func() {
		for i := 0; i < 5; i++ {
			ch <- i
		}
		close(ch)
	}()
	for v := range ch {
		fmt.Println("获取通道的值 v = ", v)
	}
}

func selectTest() {
	fmt.Println("---------select start---------")

	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(1 * time.Second)
		ch1 <- "one1"
	}()

	go func() {
		time.Sleep(2 * time.Second)
		ch2 <- "two2"
	}()

	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Println("接收 ch1:", msg1)
		case msg2 := <-ch2:
			fmt.Println("接收 ch2:", msg2)
		}
	}
	fmt.Println("---------select end---------")
}

func divideError(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("divide by zero")
	}
	return a / b, nil
}

//panic 和 recover

func panicTest() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Recovered from panic  恢复:", err)
		}
	}()

	panic("something went wrong  发生严重错误")
	fmt.Println("after panic 不会执行")

}
