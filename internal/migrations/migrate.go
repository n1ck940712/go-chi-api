package migrations

import (
	"fmt"

	"go-chi-api/internal/database"
	"go-chi-api/internal/models"
)

func init() {
	database.Connect()
}

func Migrate() {
	fmt.Println("Running migration...")
	database.DB.AutoMigrate(&models.ItemTable{})
	database.DB.AutoMigrate(&models.ItemTypeTable{})
	database.DB.AutoMigrate(&models.ItemPriceHistoryTable{})
	database.DB.AutoMigrate(&models.User{})
	database.DB.AutoMigrate(&models.SessionTable{})
	// database.DB.AutoMigrate(&models.CustomerTable{})
	// database.DB.AutoMigrate(&models.OrderTable{})
	// database.DB.AutoMigrate(&models.OrderItemTable{})
	database.DB.AutoMigrate(&models.RestockTable{})
	database.DB.AutoMigrate(&models.RestockItemTable{})
	fmt.Println("Migration complete")
}
