package servers

import (
	"04blog/constant"
	"04blog/models"
	"04blog/repositories"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type CommentService interface {
	// CreateComment 创建评论
	CreateComment(comment *models.Comment) error
	// GetCommentByID 根据评论ID获取评论信息
	GetCommentByID(commentID int64) (*models.Comment, error)
	// GetCommentsByPostID 根据文章ID获取评论列表
	GetCommentsByPostID(postID int64) ([]*models.Comment, error)
	// UpdateComment 更新评论
	UpdateComment(comment *models.Comment) error
	// DeleteComment 删除评论
	DeleteComment(commentID int64) error
}

type CommentServiceImpl struct {
	commentDao  repositories.CommentRepository
	redisClient *redis.Client
}

func NewCommentService(commentDao repositories.CommentRepository, redisClient *redis.Client) CommentService {
	return &CommentServiceImpl{commentDao: commentDao, redisClient: redisClient}
}

// CreateComment 创建评论
func (c *CommentServiceImpl) CreateComment(comment *models.Comment) error {
	// 新增评论后，需要更新文章的评论数
	if comment.PostID > 0 {
		// 使用PostID构建正确的缓存键
		c.redisClient.Incr(c.redisClient.Context(), fmt.Sprintf(constant.RedisKeyPostComments, comment.PostID))
	}

	return c.commentDao.CreateComment(comment)
}

// GetCommentsByPostID 根据文章ID获取评论列表
func (c *CommentServiceImpl) GetCommentsByPostID(postID int64) ([]*models.Comment, error) {
	return c.commentDao.GetCommentsByPostID(postID)
}

// GetCommentByID 根据评论ID获取评论信息
func (c *CommentServiceImpl) GetCommentByID(commentID int64) (*models.Comment, error) {
	return c.commentDao.GetCommentByID(commentID)
}

// UpdateComment 更新评论
func (c *CommentServiceImpl) UpdateComment(comment *models.Comment) error {
	return c.commentDao.UpdateComment(comment)
}

// DeleteComment 删除评论
func (c *CommentServiceImpl) DeleteComment(commentID int64) error {
	// 先获取评论信息，以便获取PostID
	comment, err := c.commentDao.GetCommentByID(commentID)
	if err == nil && comment != nil && comment.PostID > 0 {
		// 使用PostID构建正确的缓存键
		c.redisClient.Decr(c.redisClient.Context(), fmt.Sprintf(constant.RedisKeyPostComments, comment.PostID))
	}

	return c.commentDao.DeleteComment(commentID)
}
