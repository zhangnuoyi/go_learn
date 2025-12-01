package servers

import (
	"04blog/models"
	"04blog/repositories"
)

type PostService interface {
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
type PostServiceImpl struct {
	repo repositories.PostRepository
}

// DeletePost 删除文章
func (p *PostServiceImpl) DeletePost(postID int64) error {
	return p.repo.DeletePost(postID)
}

// GetPost 获取文章
func (p *PostServiceImpl) GetPost(postID int64) (*models.Post, error) {
	return p.repo.GetPost(postID)
}

// GetPostsByConditions 获取多条件查询文章
func (p *PostServiceImpl) GetPostsByConditions(conditions map[string]interface{}) ([]*models.Post, error) {
	return p.repo.GetPostsByConditions(conditions)
}

// UpdatePost 更新文章
func (p *PostServiceImpl) UpdatePost(post *models.Post) error {
	return p.repo.UpdatePost(post)
}

// NewPostService 创建文章服务
func NewPostService(repo repositories.PostRepository) PostService {
	return &PostServiceImpl{repo: repo}
}

// CreatePost 新增文章
func (p *PostServiceImpl) CreatePost(post *models.Post) error {
	return p.repo.CreatePost(post)
}
