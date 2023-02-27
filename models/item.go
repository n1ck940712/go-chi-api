package models

import (
	"time"
)

type ItemTable struct {
	ID        int64     `json:"id"`
	Name      string    `gorm:"text;not null;default:null" json:"name"`
	Attribute string    `gorm:"text;not null;default:null" json:"attribute"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
