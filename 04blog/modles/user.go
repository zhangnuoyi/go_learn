package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        int64          `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at"`
	IsDeleted bool           `json:"is_deleted" gorm:"column:is_deleted"`
}

// 用户相关的表数据
type User struct {
	BaseModel
	Mobile   string     `json:"mobile" gorm:"index:idx_mobile;unique; type:varchar(11);not null"`
	Password string     `json:"password" gorm:"type:varchar(100);not null"`
	Nickname string     `json:"nickname" gorm:"type:varchar(20);not null"`
	Birthday *time.Time `json:"birthday" gorm:"type:datetime"`
	Gender   string     `json:"gender" gorm:"column:gender;type:varchar(10) comment '性别:male 表示男性,female 表示女性';default:male"`
	Role     int        `json:"role" gorm:"column:role;type:int comment '角色:1 表示管理员,2 表示普通用户';default:2"`
	Address  string     `json:"address" gorm:"type:varchar(200)"`
	Email    string     `json:"email" gorm:"type:varchar(100)"`
}
