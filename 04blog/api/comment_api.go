package api

import (
	"04blog/models"
	"04blog/servers"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentVo struct {
	ID int64 `json:"id"`
	//必填 长度 2-2000
	Content string `json:"content" binding:"required,min=2,max=2000"`
	//必填
	UserID int64 `json:"user_id" binding:"required"`
	//必填
	PostID int64 `json:"post_id" binding:"required"`
}

type CommentApi struct {
	commentService servers.CommentService
}

// NewCommentApi 创建评论API
func NewCommentAPI(commentService servers.CommentService) *CommentApi {
	return &CommentApi{commentService: commentService}
}

// 新增评论
func (comm *CommentApi) CreateComment(c *gin.Context) {
	var commentVo CommentVo
	if err := c.ShouldBindJSON(&commentVo); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	//上下文获取用户id
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "该用户没有登录 没有权限"})
		return
	}
	commentVo.UserID = userID.(int64)
	//todo 需要更新评论表中的评论数使用rides实现 评论数+1 默认为0

	err := comm.commentService.CreateComment(&models.Comment{
		Content: commentVo.Content,
		UserID:  commentVo.UserID,
		PostID:  commentVo.PostID,
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "评论创建成功"})
}

// 根据文章ID获取评论
func (comm *CommentApi) GetCommentsByPostID(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("postID"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "无效的文章ID"})
		return
	}
	comments, err := comm.commentService.GetCommentsByPostID(postID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": comments})
}

// 删除评论
func (comm *CommentApi) DeleteComment(c *gin.Context) {
	commentID, err := strconv.ParseInt(c.Param("commentID"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "无效的评论ID"})
		return
	}
	err = comm.commentService.DeleteComment(commentID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "评论删除成功"})
}

// 根据评论ID获取评论
func (comm *CommentApi) GetCommentByID(c *gin.Context) {
	commentID, err := strconv.ParseInt(c.Param("commentID"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "无效的评论ID"})
		return
	}
	comment, err := comm.commentService.GetCommentByID(commentID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": comment})
}
