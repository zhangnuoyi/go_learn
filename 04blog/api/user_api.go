package api

import (
	"04blog/models"
	"04blog/servers"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserAPI struct {
	userService servers.UserService
}

func NewUserAPI(userService servers.UserService) *UserAPI {
	return &UserAPI{userService: userService}
}

// 用户注册VO
type RegisterUserVO struct {
	// 用户名 必填  长度在（3-20）
	Username string `json:"username" binding:"required,min=3,max=20"`
	// 手机号 选填  长度在（11）
	Mobile string `json:"mobile" binding:"omitempty,len=11"`
	// 密码 必填  长度在（6-20）
	Password string `json:"password" binding:"required,min=6,max=20"`
	// 昵称 选填  长度在（3-20）
	Nickname string `json:"nickname" binding:"omitempty,min=3,max=20"`
}

// 登录VO
type LoginUserVO struct {
	// 用户名或手机号 必填
	UsernameOrMobile string `json:"username_or_mobile" binding:"required"`
	// 密码 必填  长度在（6-20）
	Password string `json:"password" binding:"required,min=6,max=20"`
}

// 登录响应VO
type LoginUserResponseVO struct {
	// 登录凭证
	Token string `json:"token"`
	// 用户信息
	UserInfo models.User `json:"user_info"`
}

// 用户注册
func (u *UserAPI) RegisterUser(c *gin.Context) {
	var vo RegisterUserVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//创建用户模型
	user := models.User{
		Username: vo.Username,
		Mobile:   vo.Mobile,
		Password: vo.Password,
		Nickname: vo.Nickname,
	}
	//调用服务层注册用户
	if err := u.userService.RegisterUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//返回注册成功响应
	c.JSON(http.StatusOK, gin.H{"message": "用户注册成功"})
}

// 用户登录
func (u *UserAPI) LoginUser(c *gin.Context) {
	var vo LoginUserVO
	// 绑定JSON请求体到VO
	if err := c.ShouldBindJSON(&vo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//调用服务层登录用户
	token, userInfo, err := u.userService.LoginUser(vo.UsernameOrMobile, vo.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	//返回登录成功响应
	c.JSON(http.StatusOK, gin.H{
		"token":     token,
		"user_info": userInfo,
	})
}

//获取当前用户信息

func (userApi *UserAPI) GetCurrentUser(c *gin.Context) {
	//请求参数周获取user_id
	userIDStr := c.Query("user_id")
	// if userIDStr == {
	// 	//从上下文获取用户ID
	// 	userID, exists := c.Get("user_id")
	// 	if !exists {
	// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
	// 		return
	// 	}
	// }

	// 转换为int64类型
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}
	//调用服务层获取用户信息
	userInfo, err := userApi.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//返回用户信息
	c.JSON(http.StatusOK, gin.H{"user_info": userInfo})
}
