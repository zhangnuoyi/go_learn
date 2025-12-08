package models

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        int64          `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at"`
	IsDeleted bool           `json:"is_deleted" gorm:"column:is_deleted"`
}

// Post 博客文章模型
type Post struct {
	BaseModel
	Title     string `json:"title" gorm:"type:varchar(200) comment '文章标题';not null"`
	Content   string `json:"content" gorm:"type:text comment '文章内容';not null"`
	Summary   string `json:"summary" gorm:"type:varchar(500) comment '文章摘要'"`
	Cover     string `json:"cover" gorm:"type:varchar(255) comment '文章封面'"`
	ViewCount int    `json:"view_count" gorm:"default:0;comment '文章阅读量'"`
	LikeCount int    `json:"like_count" gorm:"default:0;comment '文章点赞量'"`
	UserID    int64  `json:"user_id" gorm:"index;not null;comment '文章作者'"`
}
