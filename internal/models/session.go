package models

import (
	"errors"
	"go-chi-api/internal/database"
	"time"

	"gorm.io/gorm"
)

type SesssionTable struct {
	ID      int32     `json:"id"`
	UserID  int32     `json:"-"`
	User    User      `gorm:"foreignKey:UserID;references:ID" json:"user"`
	Token   string    `gorm:"uniqueIndex"`
	Expires time.Time `json:"expires"`
}

// get session entry by token
func GetSessionByToken(token string) (*SesssionTable, error) {
	s := &SesssionTable{}
	err := database.DB.Model(&SesssionTable{}).Where("token = ?", token).Where("expires > ?", time.Now()).Take(&s).Error
	if err != nil {
		return &SesssionTable{}, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &SesssionTable{}, errors.New("session not found")
	}
	return s, err
}
