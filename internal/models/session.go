package models

import (
	"errors"
	"go-chi-api/internal/database"
	"time"

	"gorm.io/gorm"
)

type SessionTable struct {
	ID      int32     `json:"id"`
	UserID  int32     `json:"user_id"`
	User    User      `gorm:"foreignKey:UserID;references:ID" json:"-"`
	Token   string    `gorm:"uniqueIndex"`
	Expires time.Time `json:"expires"`
}

// get session entry by token
func GetSessionByToken(token string) (*SessionTable, error) {
	s := &SessionTable{}
	err := database.DB.Model(&SessionTable{}).Where("token = ?", token).Where("expires > ?", time.Now()).Take(&s).Error
	if err != nil {
		return &SessionTable{}, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &SessionTable{}, errors.New("session not found")
	}
	return s, err
}
