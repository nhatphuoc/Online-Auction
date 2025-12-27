package scripts

import (
	"category_service/internal/handlers"
	"category_service/internal/models"
	"context"
	"log"

	"github.com/go-pg/pg/v10"
)

func SeedInitialData(db *pg.DB) error {
	ctx := context.Background()

	log.Println("Starting seed data...")

	// Seed categories
	categories := []*models.Category{
		// Level 1 categories
		{
			Name:         "Điện tử",
			Slug:         "dien-tu",
			Description:  "Các sản phẩm điện tử",
			Level:        1,
			IsActive:     true,
			DisplayOrder: 1,
			CreatedAt:    handlers.FixedTimeNow(),
			UpdatedAt:    handlers.FixedTimeNow(),
		},
		{
			Name:         "Thời trang",
			Slug:         "thoi-trang",
			Description:  "Các sản phẩm thời trang",
			Level:        1,
			IsActive:     true,
			DisplayOrder: 2,
			CreatedAt:    handlers.FixedTimeNow(),
			UpdatedAt:    handlers.FixedTimeNow(),
		},
		{
			Name:         "Gia dụng",
			Slug:         "gia-dung",
			Description:  "Các sản phẩm gia dụng",
			Level:        1,
			IsActive:     true,
			DisplayOrder: 3,
			CreatedAt:    handlers.FixedTimeNow(),
			UpdatedAt:    handlers.FixedTimeNow(),
		},
	}

	// Insert level 1 categories
	for _, cat := range categories {
		_, err := db.ModelContext(ctx, cat).Insert()
		if err != nil {
			log.Printf("Error inserting category %s: %v", cat.Name, err)
		} else {
			log.Printf("Inserted category: %s (ID: %d)", cat.Name, cat.ID)
		}
	}

	// Get inserted categories for parent_id reference
	var dientuCategory, thoitrangCategory models.Category

	err := db.ModelContext(ctx, &dientuCategory).Where("slug = ?", "dien-tu").Select()
	if err != nil {
		log.Fatalf("Error fetching Điện tử category: %v", err)
	}

	err = db.ModelContext(ctx, &thoitrangCategory).Where("slug = ?", "thoi-trang").Select()
	if err != nil {
		log.Fatalf("Error fetching Thời trang category: %v", err)
	}

	// Level 2 categories (children)
	childCategories := []*models.Category{
		// Điện tử children
		{
			Name:         "Điện thoại di động",
			Slug:         "dien-thoai-di-dong",
			Description:  "Điện thoại thông minh, điện thoại cơ bản",
			ParentID:     &dientuCategory.ID,
			Level:        2,
			IsActive:     true,
			DisplayOrder: 1,
			CreatedAt:    handlers.FixedTimeNow(),
			UpdatedAt:    handlers.FixedTimeNow(),
		},
		{
			Name:         "Máy tính xách tay",
			Slug:         "may-tinh-xach-tay",
			Description:  "Laptop, notebook",
			ParentID:     &dientuCategory.ID,
			Level:        2,
			IsActive:     true,
			DisplayOrder: 2,
			CreatedAt:    handlers.FixedTimeNow(),
			UpdatedAt:    handlers.FixedTimeNow(),
		},
		{
			Name:         "Máy tính bảng",
			Slug:         "may-tinh-bang",
			Description:  "Tablet, iPad",
			ParentID:     &dientuCategory.ID,
			Level:        2,
			IsActive:     true,
			DisplayOrder: 3,
			CreatedAt:    handlers.FixedTimeNow(),
			UpdatedAt:    handlers.FixedTimeNow(),
		},
		// Thời trang children
		{
			Name:         "Giày",
			Slug:         "giay",
			Description:  "Giày thể thao, giày tây, giày sneaker",
			ParentID:     &thoitrangCategory.ID,
			Level:        2,
			IsActive:     true,
			DisplayOrder: 1,
			CreatedAt:    handlers.FixedTimeNow(),
			UpdatedAt:    handlers.FixedTimeNow(),
		},
		{
			Name:         "Đồng hồ",
			Slug:         "dong-ho",
			Description:  "Đồng hồ đeo tay nam, nữ",
			ParentID:     &thoitrangCategory.ID,
			Level:        2,
			IsActive:     true,
			DisplayOrder: 2,
			CreatedAt:    handlers.FixedTimeNow(),
			UpdatedAt:    handlers.FixedTimeNow(),
		},
		{
			Name:         "Túi xách",
			Slug:         "tui-xach",
			Description:  "Túi xách, ba lô, ví",
			ParentID:     &thoitrangCategory.ID,
			Level:        2,
			IsActive:     true,
			DisplayOrder: 3,
			CreatedAt:    handlers.FixedTimeNow(),
			UpdatedAt:    handlers.FixedTimeNow(),
		},
	}

	// Insert level 2 categories
	for _, cat := range childCategories {
		_, err := db.ModelContext(ctx, cat).Insert()
		if err != nil {
			log.Printf("Error inserting category %s: %v", cat.Name, err)
		} else {
			log.Printf("Inserted child category: %s (ID: %d, Parent: %d)", cat.Name, cat.ID, *cat.ParentID)
		}
	}

	log.Println("Seed data completed successfully!")
	log.Println("\nCategory structure:")
	log.Println("├── Điện tử")
	log.Println("│   ├── Điện thoại di động")
	log.Println("│   ├── Máy tính xách tay")
	log.Println("│   └── Máy tính bảng")
	log.Println("├── Thời trang")
	log.Println("│   ├── Giày")
	log.Println("│   ├── Đồng hồ")
	log.Println("│   └── Túi xách")
	log.Println("└── Gia dụng")

	return nil
}
