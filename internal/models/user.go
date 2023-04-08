package models

import (
	"time"
)

type User struct {
	ID        uint32     `json:"id"`
	Points    float32    `gorm:"type:FLOAT UNSIGNED" json:"points"`
	UpdatedAt *time.Time `json:"updated_at"`
	CreatedAt time.Time  `json:"created_at"`

	Activities []Activity `gorm:"foreignKey:UserID;references:ID;" json:"activities"`
}
