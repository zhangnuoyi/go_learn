package repositories

import (
	"04blog/models"

	"gorm.io/gorm"
)

// 定义评论的仓库接口
type CommentRepository interface {
	//新增评论
	CreateComment(comment *models.Comment) error
	//根据文章ID获取评论
	GetCommentsByPostID(postID int64) ([]*models.Comment, error)
	//删除评论
	DeleteComment(commentID int64) error
	//根据评论ID获取评论
	GetCommentByID(commentID int64) (*models.Comment, error)
	//更新评论
	UpdateComment(comment *models.Comment) error
}

// CommentRepositoryImpl 评论仓库实现
type CommentRepositoryImpl struct {
	db *gorm.DB
}

// NewCommentRepository 创建评论仓库实现
func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &CommentRepositoryImpl{db: db}
}

// 新增评论
func (c *CommentRepositoryImpl) CreateComment(comment *models.Comment) error {
	return c.db.Create(comment).Error
}

// 根据文章ID获取评论
func (c *CommentRepositoryImpl) GetCommentsByPostID(postID int64) ([]*models.Comment, error) {
	var comments []*models.Comment
	err := c.db.Where("post_id = ?", postID).Find(&comments).Error
	return comments, err
}

// 删除评论
func (c *CommentRepositoryImpl) DeleteComment(commentID int64) error {
	return c.db.Delete(&models.Comment{}, commentID).Error
}

// 根据评论ID获取评论
func (c *CommentRepositoryImpl) GetCommentByID(commentID int64) (*models.Comment, error) {
	var comment models.Comment
	err := c.db.Where("id = ?", commentID).First(&comment).Error
	return &comment, err
}

// 更新评论
func (c *CommentRepositoryImpl) UpdateComment(comment *models.Comment) error {
	return c.db.Save(comment).Error
}
