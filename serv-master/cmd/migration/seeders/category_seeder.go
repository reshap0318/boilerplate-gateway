package seeders

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/reshap0318/serv-master/internal/models"
)

// SeedCategories inserts example categories — the reference seeder to copy for new features.
// Idempotent: safe to run multiple times, existing rows (matched by name) are skipped.
func SeedCategories(db *gorm.DB) {
	fmt.Println("  Seeding categories...")

	electronicsDesc := "Electronic devices and accessories"
	booksDesc := "Books, magazines, and printed media"

	categories := []models.Category{
		{Name: "Electronics", Description: &electronicsDesc},
		{Name: "Books", Description: &booksDesc},
		{Name: "Stationery", Description: nil},
	}

	for _, category := range categories {
		var existing models.Category
		err := db.Where("name = ?", category.Name).First(&existing).Error

		switch {
		case err == nil:
			fmt.Printf("    Skipped (exists): %s\n", category.Name)
		case err == gorm.ErrRecordNotFound:
			if err := db.Create(&category).Error; err != nil {
				fmt.Printf("    Failed to seed %s: %v\n", category.Name, err)
				continue
			}
			fmt.Printf("    Created: %s\n", category.Name)
		default:
			fmt.Printf("    Failed to check %s: %v\n", category.Name, err)
		}
	}
}
