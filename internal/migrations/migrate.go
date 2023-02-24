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
	database.DB.AutoMigrate(&models.ItemTable{})
	database.DB.AutoMigrate(&models.User{})
	fmt.Println("Migration complete")
}
