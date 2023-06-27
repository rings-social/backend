package models

import (
	"time"
)

type Model struct {
	ID        uint       `json:"id" gorm:"primarykey;autoIncrmeent:true"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty" gorm:"index"`
}
