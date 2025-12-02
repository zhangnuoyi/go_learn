package servers

import (
	"04blog/constant"
	"04blog/models"
	"04blog/repositories"
	"fmt"

	"github.com/go-redis/redis/v8"
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
	repo        repositories.PostRepository
	redisClient *redis.Client
}

// DeletePost 删除文章
func (p *PostServiceImpl) DeletePost(postID int64) error {
	return p.repo.DeletePost(postID)
}

// GetPost 获取文章
func (p *PostServiceImpl) GetPost(postID int64) (*models.Post, error) {
	//从数据库中获取文章
	post, err := p.repo.GetPost(postID)
	if err != nil {
		return nil, err
	}

	//从redis中获取文章阅读量进行+1 赋值
	views, err := p.redisClient.Incr(p.redisClient.Context(), fmt.Sprintf(constant.RedisKeyPostViews, postID)).Result()
	if err != nil {
		return nil, err
	}
	post.ViewCount = int(views)
	//从redis中获取文章点赞量
	likes, err := p.redisClient.Get(p.redisClient.Context(), fmt.Sprintf(constant.RedisKeyPostLikes, postID)).Int64()
	if err != nil {
		return nil, err
	}
	post.LikeCount = int(likes)

	return post, nil
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
func NewPostService(repo repositories.PostRepository, redisClient *redis.Client) PostService {
	return &PostServiceImpl{repo: repo, redisClient: redisClient}
}

// CreatePost 新增文章
func (p *PostServiceImpl) CreatePost(post *models.Post) error {
	//初始阅读量 到redis
	p.redisClient.Set(p.redisClient.Context(), fmt.Sprintf(constant.RedisKeyPostViews, post.ID), 0, 0)
	//初始点赞量 到redis
	p.redisClient.Set(p.redisClient.Context(), fmt.Sprintf(constant.RedisKeyPostLikes, post.ID), 0, 0)
	return p.repo.CreatePost(post)
}
