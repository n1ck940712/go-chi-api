package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type RestockTable struct {
	ID           int32              `json:"id"`
	UserID       int32              `json:"user_id"`
	User         User               `gorm:"foreign:UserID;references:ID" json:"-"`
	Description  string             `gorm:"type:text;null;default:null" json:"description"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
	RestockItems []RestockItemTable `gorm:"foreignkey:RestockID" json:"restock_items"`
}

type RestockItemTable struct {
	ID         int32           `json:"id"`
	RestockID  int32           `json:"-"`
	Restock    RestockTable    `gorm:"foreign:RestockID;references:ID" json:"-"`
	ItemID     int32           `json:"item_id"`
	Item       ItemTable       `gorm:"foreign:ItemID;references:ID" json:"-"`
	Quantity   int32           `gorm:"int;not null;default:0" json:"quantity"`
	UnitPrice  decimal.Decimal `gorm:"type:decimal(10,2);not null;default:0" json:"unit_price"`
	TotalPrice decimal.Decimal `gorm:"type:decimal(10,2);not null;default:0" json:"total_price"`
}
