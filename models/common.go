package models

import "time"

// CommonFields 公共字段
type CommonFields struct {
	ID        uint       `gorm:"primary_key" json:"id"` // ID
	CreatedAt time.Time  `json:"created_at"`            // 创建时间
	UpdatedAt time.Time  `json:"updated_at"`            // 最近修改时间
	DeletedAt *time.Time `json:"-"`                     // 删除时间（删除标识）
}
