package models

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        int64          `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at"`
	IsDeleted bool           `json:"is_deleted" gorm:"column:is_deleted"`
}

type User struct {
	BaseModel
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	Age         int    `json:"age"`
	Description string `json:"description"`
	Nickname    string `json:"nickname"`
	Role        string `json:"role"`
}
