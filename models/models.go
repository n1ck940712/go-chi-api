package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TestTable struct {
	gorm.Model
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Column1   string    `json:"column1" gorm:"text;not null;default:null"`
	Column2   string    `json:"column2" gorm:"text;not null;default:null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
