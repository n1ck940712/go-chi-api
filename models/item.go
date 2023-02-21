package models

import (
	"time"
)

type ItemTable struct {
	ID        int64
	Name      string `gorm:"text;not null;default:null"`
	Attribute string `gorm:"text;not null;default:null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
