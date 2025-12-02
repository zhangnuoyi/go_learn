package api

import (
	"04blog/models"
	"04blog/servers"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LikeAPI struct {
	likeService servers.LikeService
}

func NewLikeAPI(likeService servers.LikeService) *LikeAPI {
	return &LikeAPI{likeService: likeService}
}

// Create 创建点赞
func (a *LikeAPI) Create(c *gin.Context) {
	var like models.Like
	//上下文获取用户id
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}
	like.UserID = userID.(int64)
	if err := c.ShouldBindJSON(&like); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := a.likeService.Create(&like); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "新增点赞"})
}

// Delete 删除点赞
func (a *LikeAPI) Delete(c *gin.Context) {
	var like models.Like
	//上下文获取用户id
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}
	like.UserID = userID.(int64)
	if err := c.ShouldBindJSON(&like); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := a.likeService.Delete(&like); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除点赞"})
}

// GetByArticleID 获取文章点赞列表
func (a *LikeAPI) GetByArticleID(c *gin.Context) {
	var likes []*models.Like
	//请求参数中获取文章id
	postID, err := strconv.ParseInt(c.Param("postID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return

	}
	likes, err = a.likeService.GetByPostID(postID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"likes": likes})
}
