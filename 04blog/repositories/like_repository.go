package repositories

import (
	"04blog/models"

	"gorm.io/gorm"
)

// LikeRepository 点赞仓库接口
type LikeRepository interface {
	// Create 创建点赞
	Create(like *models.Like) error
	// Delete 删除点赞
	Delete(like *models.Like) error
	// GetByPostID 根据帖子ID获取点赞列表
	GetByPostID(postID int64) ([]*models.Like, error)
	// GetByUserIDAndPostID 根据用户ID和帖子ID获取点赞
	GetByUserIDAndPostID(userID, postID int64) (*models.Like, error)
}

type LikeRepositoryImpl struct {
	db *gorm.DB
}

// NewLikeRepositoryImpl 创建点赞仓库实现
func NewLikeRepositoryImpl(db *gorm.DB) *LikeRepositoryImpl {
	return &LikeRepositoryImpl{db: db}
}

// Create 点赞
func (r *LikeRepositoryImpl) Create(like *models.Like) error {
	return r.db.Create(like).Error
}

// Delete 删除点赞
func (r *LikeRepositoryImpl) Delete(like *models.Like) error {
	return r.db.Where("user_id = ? AND post_id = ?", like.UserID, like.PostID).Delete(&models.Like{}).Error
}

// GetByPostID 根据帖子ID获取点赞列表
func (r *LikeRepositoryImpl) GetByPostID(postID int64) ([]*models.Like, error) {
	var likes []*models.Like
	err := r.db.Where("post_id = ?", postID).Find(&likes).Error
	return likes, err
}

// GetByUserIDAndPostID 根据用户ID和帖子ID获取点赞
func (r *LikeRepositoryImpl) GetByUserIDAndPostID(userID, postID int64) (*models.Like, error) {
	var like models.Like
	err := r.db.Where("user_id = ? AND post_id = ?", userID, postID).First(&like).Error
	return &like, err
}
