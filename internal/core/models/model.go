package models

import (
	"time"

	"gorm.io/gorm"
)

type Default struct {
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

type Pagination struct {
	Page      int   `json:"page" query:"page" gorm:"-"`
	Limit     int   `json:"limit" query:"limit" gorm:"-"`
	Total     int64 `json:"total" gorm:"-"`
	TotalData int64 `json:"total_data" gorm:"-"`
}
