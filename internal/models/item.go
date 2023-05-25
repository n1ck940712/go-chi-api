package models

import (
	"go-chi-api/internal/database"
	"time"

	"github.com/shopspring/decimal"
)

type ItemTable struct {
	ID          int32           `json:"id"`
	Name        string          `gorm:"type:varchar(100);not null;default:null;uniqueIndex" json:"name"`
	Description string          `gorm:"type:text;null;default:null" json:"description"`
	Quantity    int32           `gorm:"int;not null;default:0" json:"quantity"`
	Price       decimal.Decimal `gorm:"type:decimal(10,2);not null;default:0" json:"price"`
	ItemTypeID  int32           `json:"item_type_id"`
	ItemType    ItemTypeTable   `gorm:"foreignKey:ItemTypeID;references:ID" json:"-"`
	IsActive    bool            `gorm:"type:boolean;not null;default:true" json:"is_active"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type ItemPriceHistoryTable struct {
	ID        int32     `json:"id"`
	ItemID    int32     `json:"-"`
	Item      ItemTable `gorm:"foreignKey:ItemID;references:ID" json:"item"`
	Price     float32   `gorm:"float;not null;default:0" json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ItemTypeTable struct {
	ID          int32  `json:"id"`
	Name        string `gorm:"text;not null;default:null" json:"name"`
	Description string `gorm:"text;null;default:null" json:"description"`
}

func GetItemFromID(id int32) (ItemTable, error) {
	var item ItemTable
	result := database.DB.First(&item, id)
	if result.Error != nil {
		return item, result.Error
	}
	return item, nil
}
