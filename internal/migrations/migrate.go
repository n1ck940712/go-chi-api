package migrations

import (
	"fmt"

	"go-chi-api/internal/database"
	"go-chi-api/models"
)

func init() {
	database.Connect()
}

func Migrate() {
	fmt.Println("Running migration...")
	database.DB.DB.AutoMigrate(&models.ItemTable{})
	fmt.Println("Migration complete")
}
