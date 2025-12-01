package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// main 函数是程序的入口点
func main() {
	// 以下是被注释掉的标准库HTTP服务器示例
	// fmt.Println("Starting simple HTTP server...")

	// http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
	// 	name := r.URL.Query().Get("name")
	// 	if name == "" {
	// 		name = "World"
	// 	}
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(http.StatusOK)
	// 	fmt.Fprintf(w, `{"message": "Hello, %s!"}`, name)
	// })

	// fmt.Println("Server is running on http://localhost:8080")
	// if err := http.ListenAndServe(":8080", nil); err != nil {
	// 	fmt.Printf("Error starting server: %v\n", err)
	// }

	// 调用不同的示例函数来演示Gin框架的各种功能
	// helloWeb() - 创建基本的HTTP服务器
	// AsciiJSON() - 演示AsciiJSON响应功能
	// requestBody() - 演示请求体绑定功能
	// formBody() - 演示表单处理功能
	bindQueryAndForn() // 演示查询字符串和表单数据绑定
}

// helloWeb 创建一个简单的Web服务器，演示基本的GET请求处理
func helloWeb() {
	// 创建带默认中间件（日志与恢复）的Gin路由器
	r := gin.Default()

	// 定义简单的GET路由，处理/hello路径的请求
	r.GET("/hello", func(c *gin.Context) {
		// 从查询参数中获取name，如果不存在则使用默认值"World"
		name := c.Query("name")
		if name == "" {
			name = "World"
		}

		// 打印获取到的参数
		fmt.Println("参数名称 name:", name)

		// 返回JSON格式的响应
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, " + name + "!",
		})
	})

	// 启动服务器，监听8081端口
	fmt.Println("Server is running on http://localhost:8081")
	r.Run(":8081")
}

func helloWeb() {
	r := gin.Default()
	// 定义简单的 GET 路由
	r.GET("/hello", func(c *gin.Context) {
		name := c.Query("name")
		if name == "" {
			name = "World"
		}
		fmt.Println("参数名称 name:", name)
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, " + name + "!",
		})
	})
	// 启动服务器
	fmt.Println("Server is running on http://localhost:8081")
	r.Run(":8081")
}

// AsciiJSON 演示如何使用AsciiJSON生成具有转义的非ASCII字符的ASCII-only JSON
func AsciiJSON() {
	// 创建带默认中间件的Gin路由器
	r := gin.Default()

	// 定义GET路由，处理/someJosn路径的请求
	r.GET("/someJosn", func(c *gin.Context) {
		// 创建包含非ASCII字符的数据
		data := map[string]interface{}{
			"lang": "Go语言",
			"tag":  "<br>",
		}

		// 使用AsciiJSON生成具有转义的非ASCII字符的ASCII-only JSON
		c.AsciiJSON(http.StatusOK, data)
	})

	// 启动服务器，监听8080端口
	r.Run(":8080")
}

// requestBody 演示如何将请求体绑定到不同的结构体中
func requestBody() {
	// 创建带默认中间件的Gin路由器
	router := gin.Default()

	// 定义三个GET路由，分别调用不同的处理函数
	router.GET("/getb", GetDataB)
	router.GET("/getc", GetDataC)
	router.GET("/getd", GetDataD)

	// 启动服务器，使用默认端口
	router.Run()
}

// formBody 演示如何处理表单数据
func formBody() {
	// 创建带默认中间件的Gin路由器
	r := gin.Default()

	// 定义POST路由，处理根路径的请求
	r.POST("/", formHandler)

	// 启动服务器，监听8080端口
	r.Run(":8080")
}

// formA 定义表单A的数据结构
type formA struct {
	Foo string `json:"foo" xml:"foo" binding:"required"`
}

// formB 定义表单B的数据结构
type formB struct {
	Bar string `json:"bar" xml:"bar" binding:"required"`
}

// SomeHandler 演示ShouldBind的使用，注意c.Request.Body不可重用
func SomeHandler(c *gin.Context) {
	objA := formA{}
	objB := formB{}

	// 尝试绑定到formA
	// 注意：c.ShouldBind会消费c.Request.Body，一旦使用就不可重用
	if err := c.ShouldBind(&objA); err != nil {
		c.String(http.StatusOK, "the body shold be formA")
	} else if err := c.ShouldBind(&objB); err != nil {
		// 由于c.Request.Body已经被消费，这里的绑定总是会失败
		c.String(http.StatusOK, "the body shold be formB")
	}
}

// SomeHandler2 演示ShouldBindBodyWith的使用，它会缓存请求体以便重用
func SomeHandler2(c *gin.Context) {
	objA := formA{}
	objB := formB{}

	// 尝试绑定到formA，ShouldBindBodyWith会缓存请求体
	if err := c.ShouldBindBodyWith(&objA, binding.JSON); err != nil {
		c.String(http.StatusOK, "the body shold be formA")
	} else if err := c.ShouldBindBodyWith(&objB, binding.JSON); err != nil {
		// 由于请求体已被缓存，可以再次使用
		c.String(http.StatusOK, "the body shold be formB")
	} else {
		c.String(http.StatusOK, "the body is formA and formB")
	}
}

// StructA 定义简单的结构体A
type StructA struct {
	FiledA string `form:"filed_a"`
}

// StructB 定义包含嵌套结构体的结构体B
type StructB struct {
	NestedStruct StructA // 嵌套结构体
	FiledaB      string  `form:"filed_b"`
}

// StructC 定义包含指针类型嵌套结构体的结构体C
type StructC struct {
	NestedStructPointer *StructA `form:"filed_a"` // 嵌套结构体指针
	FiledC              string   `form:"filed_c"`
}

// StructD 定义包含匿名嵌套结构体的结构体D
type StructD struct {
	NestedAnonyStruct struct {
		FiledX string `form:"filed_x"`
	} // 匿名嵌套结构体
	FiledD string `form:"filed_d"`
}

// GetDataB 处理StructB类型的数据绑定
func GetDataB(c *gin.Context) {
	var b StructB

	// 绑定请求数据到结构体
	c.Bind(&b)

	// 返回JSON响应，包含绑定的数据
	c.JSON(http.StatusOK, gin.H{
		"a": b.NestedStruct, // 嵌套结构体数据
		"b": b.FiledaB,      // 普通字段数据
	})
}

// GetDataC 处理StructC类型的数据绑定
func GetDataC(con *gin.Context) {
	var c StructC

	// 绑定请求数据到结构体
	con.Bind(&c)

	// 返回JSON响应，包含绑定的数据
	con.JSON(http.StatusOK, gin.H{
		"a": c.NestedStructPointer, // 嵌套结构体指针数据
		"c": c.FiledC,              // 普通字段数据
	})
}

// GetDataD 处理StructD类型的数据绑定
func GetDataD(c *gin.Context) {
	var d StructD

	// 绑定请求数据到结构体
	c.Bind(&d)

	// 返回JSON响应，包含绑定的数据
	c.JSON(http.StatusOK, gin.H{
		"a": d.NestedAnonyStruct.FiledX, // 匿名嵌套结构体字段数据
		"d": d.FiledD,                   // 普通字段数据
	})
}

// myForm 定义用于处理数组类型表单数据的结构体
type myForm struct {
	Colors []string `form:"colors[]"` // 处理数组类型的表单字段
}

// formHandler 处理表单数据，特别是数组类型的字段
func formHandler(c *gin.Context) {
	var form myForm

	// 绑定表单数据到结构体
	c.Bind(&form)

	// 返回JSON响应，包含绑定的数组数据
	c.JSON(http.StatusOK, gin.H{
		"colors": form.Colors,
	})
}

// Persion 定义用于绑定查询字符串或表单数据的结构体
type Persion struct {
	Name      string    `form:"name"`
	Address   string    `form:"address"`
	BirthDate time.Time `form:"birth_date" time_format:"2006-01-02"` // 带时间格式的字段
}

// bindQueryAndForn 演示如何绑定查询字符串或表单数据
func bindQueryAndForn() {
	// 创建带默认中间件的Gin路由器
	r := gin.Default()

	// 定义GET路由，处理/persion路径的请求
	r.GET("/persion", persion)

	// 启动服务器，监听8080端口
	r.Run(":8080")
}

// persion 处理个人信息的查询字符串或表单数据绑定
func persion(c *gin.Context) {
	var p Persion

	// 根据请求类型自动选择绑定引擎
	// - 如果是GET请求，只使用Form绑定引擎（query）
	// - 如果是POST请求，首先检查content-type是否为JSON或XML，然后再使用Form（form-data）
	if c.ShouldBind(&p) == nil {
		// 打印绑定成功的数据
		log.Printf("name: %s, address: %s, birthDate: %s", p.Name, p.Address, p.BirthDate)
	}

	// 返回简单的成功响应
	c.String(200, "success")
}

// 以下是被注释掉的静态资源嵌入示例
// staticResource 演示如何嵌入和提供静态资源
// func staticResource() {
// 	r := gin.Default()
// 	templ := template.Must(template.New("").ParseFS(f, "templates/*.tmpl", "templates/foo/*.tmpl"))
// 	r.SetHTMLTemplate(templ)
// 	r.StaticFS("/public", http.FS(f))
// 	r.Static("/static", "./static")
// 	r.Run(":8080")
// }
