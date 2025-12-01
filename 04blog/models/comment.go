package models

// Comment 评论模型
type Comment struct {
	BaseModel
	Content string `json:"content" gorm:"type:text comment '评论内容';not null"`
	UserID  int64  `json:"user_id" gorm:"index;not null;comment '评论用户ID'"` // 评论用户ID
	PostID  int64  `json:"post_id" gorm:"index;not null;comment '评论文章ID'"` // 评论文章ID

	// 关联
	// User User `json:"user" gorm:"foreignKey:UserID"` // 评论用户
	// Post Post `json:"-" gorm:"foreignKey:PostID"`    // 所属文章
}

// TableName 指定表名
func (Comment) TableName() string {
	return "comments"
}
