package models

import (
	"gorm.io/gorm"
)

type CustomerTable struct {
	gorm.Model
	Name          string `gorm:"type:varchar(100);not null;default:null;uniqueIndex" json:"name" validate:"required,min=3,max=100"`
	Email         string `gorm:"type:varchar(100);default:null" json:"email" validate:"omitempty,email"`
	ContactNumber string `gorm:"type:varchar(100);default:null" json:"contact_number" validate:"omitempty,min=3,max=100"`
	Address       string `gorm:"type:text;default:null" json:"address" validate:"omitempty,min=3,max=500"`
	IsActive      bool   `gorm:"type:boolean;not null;default:true" json:"is_active" `
}
