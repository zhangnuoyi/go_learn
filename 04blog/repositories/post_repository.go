package repositories

import (
	"04blog/models"

	"gorm.io/gorm"
)

// 文章修改的仓库接口
type PostRepository interface {
	// 新增文章
	CreatePost(post *models.Post) error
	// 获取文章
	GetPost(postID int64) (*models.Post, error)
	//多条件查询文章 不传入就是查询所有
	GetPostsByConditions(conditions map[string]interface{}) ([]*models.Post, error)
	// 删除文章
	DeletePost(postID int64) error
	// 更新文章
	UpdatePost(post *models.Post) error
}

type PostRepositoryImpl struct {
	db *gorm.DB
}

// NewPostRepository 创建文章数据访问实例
func NewPostRepository(db *gorm.DB) PostRepository {
	return &PostRepositoryImpl{db: db}
}

// CreatePost 新增文章 常量
func (p *PostRepositoryImpl) CreatePost(post *models.Post) error {
	return p.db.Create(post).Error
}

// GetPost 获取文章
func (p *PostRepositoryImpl) GetPost(postID int64) (*models.Post, error) {
	var post models.Post
	if err := p.db.Where("id = ?", postID).First(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

// GetPostsByConditions 多条件查询文章
func (p *PostRepositoryImpl) GetPostsByConditions(conditions map[string]interface{}) ([]*models.Post, error) {
	var posts []*models.Post
	//内容模糊查询
	if content, ok := conditions["content"]; ok {
		conditions["content"] = gorm.Expr("content LIKE ?", "%"+content.(string)+"%")
	}
	if err := p.db.Where(conditions).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// UpdatePost 更新文章
func (p *PostRepositoryImpl) UpdatePost(post *models.Post) error {
	//根据id更新文章
	return p.db.Model(post).Updates(post).Error
}

// DeletePost 删除文章
func (p *PostRepositoryImpl) DeletePost(postID int64) error {
	return p.db.Delete(&models.Post{}, postID).Error
}
