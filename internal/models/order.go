package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type OrderTable struct {
	gorm.Model
	CustomerID int32            `json:"customer_id" validate:"required" `
	Customer   CustomerTable    `gorm:"foreign:CustomerID;references:ID" json:"-"`
	OrderTotal decimal.Decimal  `gorm:"type:decimal(10,2);not null;default:0" json:"order_total"`
	Status     string           `gorm:"type:varchar(100);not null;default:null" json:"status"`
	OrderItems []OrderItemTable `gorm:"foreignkey:OrderID" json:"order_items" validate:""`
}

type OrderItemTable struct {
	gorm.Model
	OrderID    uint            `json:"order_id"`
	Order      OrderTable      `gorm:"foreign:OrderID;references:ID" json:"-"`
	ItemID     int32           `json:"item_id"`
	Item       ItemTable       `gorm:"foreign:ItemID;references:ID" json:"-"`
	Quantity   int32           `gorm:"int;not null;default:0" json:"quantity"`
	UnitPrice  decimal.Decimal `gorm:"type:decimal(10,2);not null;default:0" json:"unit_price"`
	TotalPrice decimal.Decimal `gorm:"type:decimal(10,2);not null;default:0" json:"total_price"`
}
