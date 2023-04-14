package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type ItemTable struct {
	ID          int64           `json:"id"`
	Name        string          `gorm:"type:varchar(100);not null;default:null" json:"name"`
	Description string          `gorm:"type:text;null;default:null" json:"description"`
	Quantity    int64           `gorm:"int;not null;default:0" json:"quantity"`
	Price       decimal.Decimal `gorm:"type:decimal(10,2);not null;default:0" json:"price"`
	ItemTypeID  int64           `json:"-"`
	ItemType    ItemTypeTable   `gorm:"foreignKey:ItemTypeID;references:ID" json:"item_type"`
	IsActive    bool            `gorm:"type:boolean;not null;default:true" json:"is_active"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type ItemPriceHistoryTable struct {
	ID        int64     `json:"id"`
	ItemID    int64     `json:"-"`
	Item      ItemTable `gorm:"foreignKey:ItemID;references:ID" json:"item"`
	Price     float32   `gorm:"float;not null;default:0" json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ItemTypeTable struct {
	ID          int64  `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Name        string `gorm:"text;not null;default:null" json:"name"`
	Description string `gorm:"text;null;default:null" json:"description"`
}
