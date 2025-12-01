package routes

import (
	"04blog/api"
	"04blog/repositories"
	"04blog/servers"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// gin routes 设置
func SetRoutes(r *gin.Engine, db *gorm.DB) {
	//后续接口路由
	//添加一个hello路由
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
	})

	//获取用户示例
	userDao := repositories.NewUserRepository(db)
	userService := servers.NewUserService(userDao)
	userApi := api.NewUserAPI(userService)

	//文章相关
	postDao := repositories.NewPostRepository(db)
	postService := servers.NewPostService(postDao)
	postApi := api.NewPostAPI(postService)

	// 设置v1路由组
	v1 := r.Group("/v1")
	{
		//用户路由
		v1.POST("/register", userApi.RegisterUser)
		v1.POST("/login", userApi.LoginUser)
		v1.GET("/user", userApi.GetCurrentUser)

		//文章路由
		v1.POST("/post", postApi.SavePost)
		v1.GET("/posts", postApi.GetPosts)
		v1.GET("/post/:id", postApi.GetPost)
		v1.DELETE("/post/:id", postApi.DeletePost)
	}

}
