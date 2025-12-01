package api

import (
	"04blog/models"
	"04blog/servers"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type PostAPI struct {
	service servers.PostService
}
type PostVO struct {
	ID int64 `json:"id"`
	//必填 长度 2-50
	Title string `json:"title" binding:"required,min=2,max=50"`
	//必填 长度 2-2000
	Content string `json:"content" binding:"required,min=2,max=2000"`
	//必填 长度 2-500
	Summary string `json:"summary" binding:"required,min=2,max=500"`
	//可选 长度 2-255
	Cover string `json:"cover" binding:"omitempty,min=2,max=255"`

	ViewCount int   `json:"view_count"`
	LikeCount int   `json:"like_count"`
	UserID    int64 `json:"user_id" binding:"required"`
}

// NewPostAPI 创建文章API
func NewPostAPI(service servers.PostService) *PostAPI {
	return &PostAPI{service: service}
}

// SavePost 创建或更新文章
func (api *PostAPI) SavePost(c *gin.Context) {
	var vo PostVO
	//上下文中获取用户id
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "该用户没有登录 没有权限"})
		return
	}
	vo.UserID = userID.(int64)

	if err := c.ShouldBindJSON(&vo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//创建文章
	now := time.Now()
	post := models.Post{
		BaseModel: models.BaseModel{
			ID:        vo.ID,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Title:     vo.Title,
		Content:   vo.Content,
		Summary:   vo.Summary,
		Cover:     vo.Cover,
		ViewCount: vo.ViewCount,
		LikeCount: vo.LikeCount,
		UserID:    vo.UserID,
	}

	// 设置默认值
	if vo.ViewCount == 0 {
		post.ViewCount = 0
	}
	if vo.LikeCount == 0 {
		post.LikeCount = 0
	}
	message := "文章创建成功"

	if vo.ID == 0 {
		//创建文章
		fmt.Println("创建文章", post)
		if err := api.service.CreatePost(&post); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		//更新文章
		post.UpdatedAt = time.Now() // 确保更新时也设置正确的更新时间
		fmt.Println("更新文章", post)
		message = "文章更新成功"
		if err := api.service.UpdatePost(&post); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": message})
}

// GetPost 获取文章
func (api *PostAPI) GetPost(c *gin.Context) {
	//获取路径动态参数id
	postIDStr := c.Param("id")

	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	post, err := api.service.GetPost(postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, post)
}

// 多条件查询文章 不传入就是查询所有
func (api *PostAPI) GetPosts(c *gin.Context) {
	conditions := map[string]interface{}{}
	//请求参数  可选参数
	title := c.Query("title")
	if title != "" {
		conditions["title"] = title
	}
	//请求参数  可选参数
	userID := c.Query("user_id")
	if userID != "" {
		userIDInt, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		conditions["user_id"] = userIDInt
	}
	posts, err := api.service.GetPostsByConditions(conditions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//内容 可选参数
	content := c.Query("content")
	if content != "" {
		conditions["content"] = content
	}
	c.JSON(http.StatusOK, posts)
}

// DeletePost 删除文章
func (api *PostAPI) DeletePost(c *gin.Context) {
	id := c.Param("id")

	postID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := api.service.DeletePost(postID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "文章删除成功"})
}
