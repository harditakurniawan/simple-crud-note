package entities

import (
	"time"

	"gorm.io/gorm"
)

type Base struct {
	ID        uint           `json:"id" gorm:"primaryKey;"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime;"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime:nano;index;"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index;"`
}
