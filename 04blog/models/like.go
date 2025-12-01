package models

// Like 点赞模型
type Like struct {
	BaseModel
	UserID int64 `json:"user_id" gorm:"index;not null;comment '点赞用户ID'"` // 点赞用户ID
	PostID int64 `json:"post_id" gorm:"index;not null;comment '点赞文章ID'"` // 点赞文章ID

	// 关联
	// User User `json:"user" gorm:"foreignKey:UserID"` // 点赞用户
	// Post Post `json:"-" gorm:"foreignKey:PostID"`    // 点赞文章
}

// TableName 指定表名
func (Like) TableName() string {
	return "likes"
}
