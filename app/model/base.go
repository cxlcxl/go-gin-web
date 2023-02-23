package model

import (
	"gorm.io/gorm"
	"time"
)

type gormDb struct {
	*gorm.DB `json:"-" gorm:"-"`
}

type Timestamp struct {
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 最后更新时间
}
