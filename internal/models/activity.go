package models

import (
	"time"

	"github.com/ehabterra/rewards/internal/pb"
)

type Activity struct {
	ID        int32           `json:"id"`
	UserID    int32           `json:"user_id"`
	Type      pb.ActivityType `json:"type"`
	ObjectID  int32           `json:"object_id"`
	Points    float32         `json:"points"`
	CreatedAt time.Time       `json:"created_at"`

	User User `gorm:"foreignKey:UserID;references:ID;"`
}
