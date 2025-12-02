package servers

import (
	"04blog/constant"
	"04blog/models"
	"04blog/repositories"
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type LikeService interface {
	// Create 创建点赞
	Create(like *models.Like) error
	// Delete 删除点赞
	Delete(like *models.Like) error
	// GetByPostID 根据帖子ID获取点赞列表
	GetByPostID(postID int64) ([]*models.Like, error)
}

type LikeServiceImpl struct {
	repo        repositories.LikeRepository
	redisClient *redis.Client
}

// NewLikeServiceImpl 创建点赞服务实现
func NewLikeServiceImpl(repo repositories.LikeRepository, redisClient *redis.Client) *LikeServiceImpl {
	return &LikeServiceImpl{repo: repo, redisClient: redisClient}
}

// Create 创建点赞
func (s *LikeServiceImpl) Create(like *models.Like) error {
	// 检查是否已点赞
	existingLike, err := s.repo.GetByUserIDAndPostID(like.UserID, like.PostID)
	if err == nil && existingLike != nil {
		return errors.New("user has already liked this post")
	}
	//reids 缓存 点赞+1
	redisKey := fmt.Sprintf(constant.RedisKeyPostLikes, like.PostID)
	err = s.redisClient.Incr(context.Background(), redisKey).Err()
	if err != nil {
		return err
	}
	return s.repo.Create(like)
}

// Delete 删除点赞
func (s *LikeServiceImpl) Delete(like *models.Like) error {
	// 检查是否已点赞
	existingLike, err := s.repo.GetByUserIDAndPostID(like.UserID, like.PostID)
	if err != nil || existingLike == nil {
		return errors.New("user has not liked this post")
	}

	// 删除数据库中的点赞记录
	err = s.repo.Delete(existingLike)
	if err != nil {
		return err
	}

	// 更新Redis缓存，点赞数-1
	redisKey := fmt.Sprintf(constant.RedisKeyPostLikes, like.PostID)
	err = s.redisClient.Decr(context.Background(), redisKey).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetByPostID 根据帖子ID获取点赞列表
func (s *LikeServiceImpl) GetByPostID(postID int64) ([]*models.Like, error) {
	return s.repo.GetByPostID(postID)
}
