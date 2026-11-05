package model

import "time"

// BaseModel 通用字段。业务表可按需内嵌。
type BaseModel struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}
