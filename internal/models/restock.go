package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type RestockTable struct {
	ID          int32     `json:"id"`
	UserID      int32     `json:"-"`
	User        User      `gorm:"foreign:UserID;references:ID" json:"user"`
	Description string    `gorm:"type:text;null;default:null" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RestockItemTable struct {
	ID         int32           `json:"id"`
	RestockID  int32           `json:"-"`
	Restock    RestockTable    `gorm:"foreign:RestockID;references:ID" json:"restock"`
	ItemID     int32           `json:"-"`
	Item       ItemTable       `gorm:"foreign:ItemID;references:ID" json:"item"`
	Quantity   int32           `gorm:"int;not null;default:0" json:"quantity"`
	UnitPrice  decimal.Decimal `gorm:"type:decimal(10,2);not null;default:0" json:"unit_price"`
	TotalPrice decimal.Decimal `gorm:"type:decimal(10,2);not null;default:0" json:"total_price"`
}
