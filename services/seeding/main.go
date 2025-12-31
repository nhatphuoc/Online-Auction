package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"seeding/config"
	"time"
)

func main() {
	cfg := config.LoadConfig()
	db := config.ConnectDB(cfg)
	defer db.Close()

	ctx := context.Background()

	log.Println("Starting database seeding...")

	// Clear existing data (in reverse order of dependencies)
	if err := clearData(ctx, db); err != nil {
		log.Fatalf("Failed to clear data: %v", err)
	}

	// Seed data in order of dependencies
	if err := seedCategories(ctx, db); err != nil {
		log.Fatalf("Failed to seed categories: %v", err)
	}

	productIDMap, err := seedProducts(ctx, db)
	if err != nil {
		log.Fatalf("Failed to seed products: %v", err)
	}

	if err := seedProductImages(ctx, db, productIDMap); err != nil {
		log.Fatalf("Failed to seed product images: %v", err)
	}

	if err := seedWatchList(ctx, db); err != nil {
		log.Fatalf("Failed to seed watch list: %v", err)
	}

	if err := seedBiddingHistory(ctx, db); err != nil {
		log.Fatalf("Failed to seed bidding history: %v", err)
	}

	if err := seedComments(ctx, db); err != nil {
		log.Fatalf("Failed to seed comments: %v", err)
	}

	if err := seedOrders(ctx, db); err != nil {
		log.Fatalf("Failed to seed orders: %v", err)
	}

	if err := seedOrderMessages(ctx, db); err != nil {
		log.Fatalf("Failed to seed order messages: %v", err)
	}

	if err := seedOrderRatings(ctx, db); err != nil {
		log.Fatalf("Failed to seed order ratings: %v", err)
	}

	if err := seedUserUpgradeRequests(ctx, db); err != nil {
		log.Fatalf("Failed to seed user upgrade requests: %v", err)
	}

	log.Println("Database seeding completed successfully!")
}

func clearData(ctx context.Context, db *sql.DB) error {
	log.Println("Clearing existing data...")

	// Use TRUNCATE CASCADE to delete all data and reset sequences
	// This is faster and handles foreign key constraints automatically
	tables := []string{
		"order_ratings",
		"order_messages",
		"orders",
		"comments",
		"bidding_history",
		"watch_list",
		"product_images",
		"products",
		"categories",
		"user_upgrade_requests",
	}

	for _, table := range tables {
		_, err := db.ExecContext(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			// If TRUNCATE fails, try DELETE as fallback
			log.Printf("Warning: TRUNCATE failed for %s, trying DELETE: %v", table, err)
			_, err = db.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s", table))
			if err != nil {
				return fmt.Errorf("error clearing %s: %v", table, err)
			}
		}
		log.Printf("Cleared %s", table)
	}

	// Reset sequences
	sequences := []string{
		"categories_id_seq",
		"products_id_seq",
		"bidding_history_id_seq",
		"comments_id_seq",
		"orders_id_seq",
		"order_messages_id_seq",
		"order_ratings_id_seq",
		"watch_list_id_seq",
		"user_upgrade_requests_id_seq",
	}

	for _, seq := range sequences {
		_, err := db.ExecContext(ctx, fmt.Sprintf("ALTER SEQUENCE %s RESTART WITH 1", seq))
		if err != nil {
			// Ignore errors for sequences that might not exist
			log.Printf("Warning: Could not reset sequence %s: %v", seq, err)
		}
	}

	return nil
}

func seedCategories(ctx context.Context, db *sql.DB) error {
	log.Println("Seeding categories...")

	now := time.Now()

	// Level 1 categories (parent categories)
	level1Categories := []struct {
		name        string
		slug        string
		description string
	}{
		{"Điện tử", "dien-tu", "Các sản phẩm điện tử, thiết bị công nghệ"},
		{"Thời trang", "thoi-trang", "Quần áo, phụ kiện thời trang"},
		{"Đồ gia dụng", "do-gia-dung", "Đồ dùng trong gia đình"},
		{"Xe cộ", "xe-co", "Xe máy, ô tô và phụ kiện"},
		{"Sách & Văn phòng phẩm", "sach-van-phong-pham", "Sách, vở, dụng cụ học tập"},
		{"Thể thao & Du lịch", "the-thao-du-lich", "Đồ thể thao, dụng cụ du lịch"},
		{"Đồ chơi & Trẻ em", "do-choi-va-tre-em", "Đồ chơi, đồ dùng cho trẻ em"},
		{"Mỹ phẩm & Làm đẹp", "my-pham-lam-dep", "Mỹ phẩm, chăm sóc sắc đẹp"},
	}

	for i, cat := range level1Categories {
		_, err := db.ExecContext(ctx, `
			INSERT INTO categories (name, slug, description, parent_id, level, is_active, display_order, created_at, updated_at)
			VALUES ($1, $2, $3, NULL, 1, true, $4, $5, $6)
		`, cat.name, cat.slug, cat.description, i+1, now, now)
		if err != nil {
			return fmt.Errorf("error inserting category %s: %v", cat.name, err)
		}
	}

	// Level 2 categories (sub-categories)
	level2Categories := []struct {
		name        string
		slug        string
		description string
		parentID    int64
	}{
		// Điện tử (parent_id: 1)
		{"Điện thoại di động", "dien-thoai-di-dong", "Smartphone, điện thoại thông minh", 1},
		{"Laptop", "laptop", "Máy tính xách tay", 1},
		{"Máy tính bảng", "may-tinh-bang", "Tablet, iPad", 1},
		{"Tai nghe", "tai-nghe", "Tai nghe, headphone, earphone", 1},
		{"Smart Watch", "smart-watch", "Đồng hồ thông minh", 1},

		// Thời trang (parent_id: 2)
		{"Quần áo nam", "quan-ao-nam", "Áo, quần cho nam giới", 2},
		{"Quần áo nữ", "quan-ao-nu", "Áo, quần cho nữ giới", 2},
		{"Giày dép", "giay-dep", "Giày, dép, sandal", 2},
		{"Túi xách", "tui-xach", "Balo, túi xách, ví", 2},
		{"Đồng hồ", "dong-ho", "Đồng hồ đeo tay", 2},

		// Đồ gia dụng (parent_id: 3)
		{"Nội thất", "noi-that", "Bàn, ghế, tủ, giường", 3},
		{"Đồ điện gia dụng", "do-dien-gia-dung", "Máy giặt, tủ lạnh, lò vi sóng", 3},
		{"Đồ dùng nhà bếp", "do-dung-nha-bep", "Nồi, chảo, dao, thớt", 3},
		{"Trang trí nhà cửa", "trang-tri-nha-cua", "Tranh, đèn, đồ trang trí", 3},

		// Xe cộ (parent_id: 4)
		{"Xe máy", "xe-may", "Xe máy các loại", 4},
		{"Ô tô", "o-to", "Xe ô tô các loại", 4},
		{"Phụ tùng xe", "phu-tung-xe", "Phụ tùng, phụ kiện xe máy, ô tô", 4},

		// Sách & Văn phòng phẩm (parent_id: 5)
		{"Sách", "sach", "Sách văn học, tham khảo, giáo khoa", 5},
		{"Văn phòng phẩm", "van-phong-pham", "Bút, vở, giấy, kẹp", 5},

		// Thể thao & Du lịch (parent_id: 6)
		{"Dụng cụ thể thao", "dung-cu-the-thao", "Bóng, vợt, giày thể thao", 6},
		{"Dụng cụ du lịch", "dung-cu-du-lich", "Ba lô, lều, túi ngủ", 6},

		// Đồ chơi & Trẻ em (parent_id: 7)
		{"Đồ chơi trẻ em", "do-choi-tre-em", "Đồ chơi cho bé", 7},
		{"Đồ dùng trẻ em", "do-dung-tre-em", "Bình sữa, tã, quần áo trẻ em", 7},

		// Mỹ phẩm & Làm đẹp (parent_id: 8)
		{"Mỹ phẩm", "my-pham", "Kem dưỡng, son môi, phấn", 8},
		{"Chăm sóc da", "cham-soc-da", "Sữa rửa mặt, toner, serum", 8},
	}

	for i, cat := range level2Categories {
		_, err := db.ExecContext(ctx, `
			INSERT INTO categories (name, slug, description, parent_id, level, is_active, display_order, created_at, updated_at)
			VALUES ($1, $2, $3, $4, 2, true, $5, $6, $7)
		`, cat.name, cat.slug, cat.description, cat.parentID, i+1, now, now)
		if err != nil {
			return fmt.Errorf("error inserting category %s: %v", cat.name, err)
		}
	}

	log.Println("Categories seeded successfully")
	return nil
}

func seedProducts(ctx context.Context, db *sql.DB) (map[int]int64, error) {
	log.Println("Seeding products...")

	now := time.Now()
	productIDMap := make(map[int]int64) // index -> actual product ID

	products := []struct {
		name          string
		description   string
		categoryID    int64
		sellerID      int64
		startingPrice float64
		currentPrice  float64
		buyNowPrice   *float64
		stepPrice     float64
		status        string
		thumbnailURL  string
		autoExtend    bool
		endAt         time.Time
		currentBidder *int64
		bidCount      int64
		categoryName  string
		parentCatID   int64
		parentCatName string
	}{
		// ĐIỆN TỬ - Điện thoại di động
		{
			name:          "iPhone 15 Pro Max 256GB",
			description:   "iPhone 15 Pro Max chính hãng VN/A, còn bảo hành 11 tháng. Máy đẹp 99%, không trầy xước, đầy đủ phụ kiện.",
			categoryID:    9,
			sellerID:      10,
			startingPrice: 25000000,
			currentPrice:  27500000,
			buyNowPrice:   floatPtr(32000000),
			stepPrice:     500000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/iphone15promax.jpg",
			autoExtend:    true,
			endAt:         now.Add(48 * time.Hour),
			currentBidder: int64Ptr(7),
			bidCount:      5,
			categoryName:  "Điện thoại di động",
			parentCatID:   1,
			parentCatName: "Điện tử",
		},
		{
			name:          "Samsung Galaxy S24 Ultra 512GB",
			description:   "Samsung S24 Ultra mới 100%, fullbox nguyên seal. Màu Titanium Gray, bộ nhớ 512GB.",
			categoryID:    9,
			sellerID:      11,
			startingPrice: 24000000,
			currentPrice:  26000000,
			buyNowPrice:   floatPtr(30000000),
			stepPrice:     500000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/s24ultra.jpg",
			autoExtend:    false,
			endAt:         now.Add(36 * time.Hour),
			currentBidder: int64Ptr(8),
			bidCount:      4,
			categoryName:  "Điện thoại di động",
			parentCatID:   1,
			parentCatName: "Điện tử",
		},
		{
			name:          "iPhone 14 Pro 128GB",
			description:   "iPhone 14 Pro màu Deep Purple, đẹp như mới, đã qua sử dụng 6 tháng. Còn bảo hành 6 tháng.",
			categoryID:    9,
			sellerID:      17,
			startingPrice: 18000000,
			currentPrice:  20000000,
			buyNowPrice:   floatPtr(23000000),
			stepPrice:     500000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/iphone14pro.jpg",
			autoExtend:    true,
			endAt:         now.Add(24 * time.Hour),
			currentBidder: int64Ptr(9),
			bidCount:      4,
			categoryName:  "Điện thoại di động",
			parentCatID:   1,
			parentCatName: "Điện tử",
		},
		{
			name:          "Xiaomi 14 Ultra 16GB/512GB",
			description:   "Xiaomi 14 Ultra chính hãng, camera Leica đỉnh cao. Màu đen, bộ nhớ 16GB RAM + 512GB ROM.",
			categoryID:    9,
			sellerID:      10,
			startingPrice: 20000000,
			currentPrice:  20000000,
			buyNowPrice:   floatPtr(25000000),
			stepPrice:     500000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/xiaomi14ultra.jpg",
			autoExtend:    false,
			endAt:         now.Add(96 * time.Hour),
			currentBidder: nil,
			bidCount:      0,
			categoryName:  "Điện thoại di động",
			parentCatID:   1,
			parentCatName: "Điện tử",
		},
		// ĐIỆN TỬ - Laptop
		{
			name:          "MacBook Pro M3 14 inch 2024",
			description:   "MacBook Pro M3 chip mới nhất, RAM 16GB, SSD 512GB. Fullbox, chưa kích hoạt bảo hành.",
			categoryID:    10,
			sellerID:      10,
			startingPrice: 35000000,
			currentPrice:  38000000,
			buyNowPrice:   floatPtr(42000000),
			stepPrice:     500000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/macbook-pro-m3.jpg",
			autoExtend:    false,
			endAt:         now.Add(72 * time.Hour),
			currentBidder: int64Ptr(8),
			bidCount:      6,
			categoryName:  "Laptop",
			parentCatID:   1,
			parentCatName: "Điện tử",
		},
		// ĐIỆN TỬ - Máy tính bảng
		{
			name:          "iPad Pro M2 11 inch 128GB WiFi",
			description:   "iPad Pro M2 mới 100%, chưa active. Màu xám, WiFi, bộ nhớ 128GB.",
			categoryID:    11,
			sellerID:      11,
			startingPrice: 18000000,
			currentPrice:  20000000,
			buyNowPrice:   nil,
			stepPrice:     200000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/ipad-pro-m2.jpg",
			autoExtend:    true,
			endAt:         now.Add(24 * time.Hour),
			currentBidder: int64Ptr(7),
			bidCount:      10,
			categoryName:  "Máy tính bảng",
			parentCatID:   1,
			parentCatName: "Điện tử",
		},
		{
			name:          "iPad Air 5 M1 64GB WiFi",
			description:   "iPad Air 5 chip M1, màu Starlight, bộ nhớ 64GB. Máy đẹp 99%, còn bảo hành 7 tháng.",
			categoryID:    11,
			sellerID:      10,
			startingPrice: 12000000,
			currentPrice:  13500000,
			buyNowPrice:   floatPtr(15000000),
			stepPrice:     200000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/ipad-air5.jpg",
			autoExtend:    false,
			endAt:         now.Add(42 * time.Hour),
			currentBidder: int64Ptr(8),
			bidCount:      7,
			categoryName:  "Máy tính bảng",
			parentCatID:   1,
			parentCatName: "Điện tử",
		},
		{
			name:          "Samsung Galaxy Tab S9+ 12.4 inch",
			description:   "Samsung Tab S9+ màn hình 12.4 inch, bộ nhớ 256GB. Kèm bút S-Pen, fullbox.",
			categoryID:    11,
			sellerID:      18,
			startingPrice: 15000000,
			currentPrice:  16000000,
			buyNowPrice:   floatPtr(18000000),
			stepPrice:     200000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/tab-s9-plus.jpg",
			autoExtend:    true,
			endAt:         now.Add(54 * time.Hour),
			currentBidder: int64Ptr(9),
			bidCount:      5,
			categoryName:  "Máy tính bảng",
			parentCatID:   1,
			parentCatName: "Điện tử",
		},
		{
			name:          "ASUS ROG Zephyrus G14 2024",
			description:   "ASUS ROG G14 2024, Ryzen 9, RTX 4060, RAM 32GB, SSD 1TB. Laptop gaming siêu mỏng nhẹ.",
			categoryID:    10,
			sellerID:      18,
			startingPrice: 28000000,
			currentPrice:  30000000,
			buyNowPrice:   floatPtr(35000000),
			stepPrice:     500000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/asus-rog-g14.jpg",
			autoExtend:    false,
			endAt:         now.Add(48 * time.Hour),
			currentBidder: int64Ptr(8),
			bidCount:      4,
			categoryName:  "Laptop",
			parentCatID:   1,
			parentCatName: "Điện tử",
		},
		{
			name:          "Lenovo ThinkPad X1 Carbon Gen 11",
			description:   "ThinkPad X1 Carbon Gen 11, i7-1365U, RAM 16GB, SSD 512GB. Laptop doanh nhân cao cấp.",
			categoryID:    10,
			sellerID:      10,
			startingPrice: 25000000,
			currentPrice:  27000000,
			buyNowPrice:   nil,
			stepPrice:     500000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/thinkpad-x1.jpg",
			autoExtend:    true,
			endAt:         now.Add(84 * time.Hour),
			currentBidder: int64Ptr(9),
			bidCount:      4,
			categoryName:  "Laptop",
			parentCatID:   1,
			parentCatName: "Điện tử",
		},
		// ĐIỆN TỬ - Tai nghe
		{
			name:          "AirPods Pro 2 USB-C",
			description:   "AirPods Pro thế hệ 2 cổng USB-C, chính hãng Apple VN. Fullbox, seal nguyên.",
			categoryID:    12,
			sellerID:      10,
			startingPrice: 5000000,
			currentPrice:  5500000,
			buyNowPrice:   floatPtr(6500000),
			stepPrice:     100000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/airpods-pro-2.jpg",
			autoExtend:    false,
			endAt:         now.Add(36 * time.Hour),
			currentBidder: int64Ptr(9),
			bidCount:      5,
			categoryName:  "Tai nghe",
			parentCatID:   1,
			parentCatName: "Điện tử",
		},
		{
			name:          "Sony WH-1000XM5 Silver",
			description:   "Sony WH-1000XM5 màu bạc, tai nghe chống ồn hàng đầu. Fullbox, còn bảo hành 10 tháng.",
			categoryID:    12,
			sellerID:      11,
			startingPrice: 6000000,
			currentPrice:  6800000,
			buyNowPrice:   floatPtr(8000000),
			stepPrice:     200000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/sony-1000xm5.jpg",
			autoExtend:    true,
			endAt:         now.Add(66 * time.Hour),
			currentBidder: int64Ptr(7),
			bidCount:      4,
			categoryName:  "Tai nghe",
			parentCatID:   1,
			parentCatName: "Điện tử",
		},
		{
			name:          "Bose QuietComfort 45",
			description:   "Bose QC45 tai nghe chống ồn cao cấp. Màu đen, đã qua sử dụng 3 tháng, như mới.",
			categoryID:    12,
			sellerID:      17,
			startingPrice: 4000000,
			currentPrice:  4500000,
			buyNowPrice:   floatPtr(5500000),
			stepPrice:     100000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/bose-qc45.jpg",
			autoExtend:    false,
			endAt:         now.Add(20 * time.Hour),
			currentBidder: int64Ptr(8),
			bidCount:      5,
			categoryName:  "Tai nghe",
			parentCatID:   1,
			parentCatName: "Điện tử",
		},
		// ĐIỆN TỬ - Smart Watch
		{
			name:          "Apple Watch Series 9 GPS 45mm",
			description:   "Apple Watch Series 9 màu Midnight, dây Sport Band. Còn bảo hành 10 tháng.",
			categoryID:    13,
			sellerID:      11,
			startingPrice: 8000000,
			currentPrice:  9000000,
			buyNowPrice:   floatPtr(10500000),
			stepPrice:     200000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/apple-watch-s9.jpg",
			autoExtend:    true,
			endAt:         now.Add(60 * time.Hour),
			currentBidder: int64Ptr(8),
			bidCount:      5,
			categoryName:  "Smart Watch",
			parentCatID:   1,
			parentCatName: "Điện tử",
		},
		{
			name:          "Samsung Galaxy Watch 6 Classic 47mm",
			description:   "Galaxy Watch 6 Classic 47mm màu đen, vòng xoay cổ điển. Fullbox, mới 100%.",
			categoryID:    13,
			sellerID:      18,
			startingPrice: 6000000,
			currentPrice:  6500000,
			buyNowPrice:   floatPtr(8000000),
			stepPrice:     200000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/galaxy-watch6.jpg",
			autoExtend:    false,
			endAt:         now.Add(44 * time.Hour),
			currentBidder: int64Ptr(7),
			bidCount:      2,
			categoryName:  "Smart Watch",
			parentCatID:   1,
			parentCatName: "Điện tử",
		},
		{
			name:          "Garmin Fenix 7X Sapphire Solar",
			description:   "Garmin Fenix 7X Sapphire Solar, đồng hồ thể thao cao cấp. Pin siêu khỏe, đầy đủ tính năng.",
			categoryID:    13,
			sellerID:      10,
			startingPrice: 15000000,
			currentPrice:  16000000,
			buyNowPrice:   nil,
			stepPrice:     500000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/garmin-fenix7x.jpg",
			autoExtend:    true,
			endAt:         now.Add(90 * time.Hour),
			currentBidder: int64Ptr(8),
			bidCount:      2,
			categoryName:  "Smart Watch",
			parentCatID:   1,
			parentCatName: "Điện tử",
		},
		{
			name:          "iPad Pro M2 11 inch 128GB WiFi",
			description:   "iPad Pro M2 mới 100%, chưa active. Màu xám, WiFi, bộ nhớ 128GB.",
			categoryID:    11, // Máy tính bảng
			sellerID:      11,
			startingPrice: 18000000,
			currentPrice:  20000000,
			buyNowPrice:   nil,
			stepPrice:     200000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/ipad-pro-m2.jpg",
			autoExtend:    true,
			endAt:         now.Add(24 * time.Hour),
			currentBidder: int64Ptr(7),
			bidCount:      10,
			categoryName:  "Máy tính bảng",
			parentCatID:   1,
			parentCatName: "Điện tử",
		},
		{
			name:          "AirPods Pro 2 USB-C",
			description:   "AirPods Pro thế hệ 2 cổng USB-C, chính hãng Apple VN. Fullbox, seal nguyên.",
			categoryID:    12, // Tai nghe
			sellerID:      10,
			startingPrice: 5000000,
			currentPrice:  5500000,
			buyNowPrice:   floatPtr(6500000),
			stepPrice:     100000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/airpods-pro-2.jpg",
			autoExtend:    false,
			endAt:         now.Add(36 * time.Hour),
			currentBidder: int64Ptr(9),
			bidCount:      5,
			categoryName:  "Tai nghe",
			parentCatID:   1,
			parentCatName: "Điện tử",
		},
		{
			name:          "Apple Watch Series 9 GPS 45mm",
			description:   "Apple Watch Series 9 màu Midnight, dây Sport Band. Còn bảo hành 10 tháng.",
			categoryID:    13, // Smart Watch
			sellerID:      11,
			startingPrice: 8000000,
			currentPrice:  9000000,
			buyNowPrice:   floatPtr(10500000),
			stepPrice:     200000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/apple-watch-s9.jpg",
			autoExtend:    true,
			endAt:         now.Add(60 * time.Hour),
			currentBidder: int64Ptr(8),
			bidCount:      5,
			categoryName:  "Smart Watch",
			parentCatID:   1,
			parentCatName: "Điện tử",
		},
		// THỜI TRANG - Quần áo nam
		{
			name:          "Áo khoác Bomber Nam cao cấp",
			description:   "Áo khoác Bomber chất liệu dù cao cấp, form Hàn Quốc. Size M, L, XL.",
			categoryID:    14,
			sellerID:      10,
			startingPrice: 300000,
			currentPrice:  450000,
			buyNowPrice:   floatPtr(600000),
			stepPrice:     50000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/bomber-jacket.jpg",
			autoExtend:    false,
			endAt:         now.Add(12 * time.Hour),
			currentBidder: int64Ptr(7),
			bidCount:      3,
			categoryName:  "Quần áo nam",
			parentCatID:   2,
			parentCatName: "Thời trang",
		},
		{
			name:          "Áo Polo nam Ralph Lauren",
			description:   "Áo Polo Ralph Lauren chính hãng Mỹ, chất liệu cotton cao cấp. Size L, màu trắng.",
			categoryID:    14,
			sellerID:      11,
			startingPrice: 400000,
			currentPrice:  550000,
			buyNowPrice:   floatPtr(700000),
			stepPrice:     50000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/polo-ralph.jpg",
			autoExtend:    false,
			endAt:         now.Add(28 * time.Hour),
			currentBidder: int64Ptr(8),
			bidCount:      3,
			categoryName:  "Quần áo nam",
			parentCatID:   2,
			parentCatName: "Thời trang",
		},
		{
			name:          "Quần Jean nam Levi's 511",
			description:   "Quần jean Levi's 511 Slim Fit, màu xanh đậm. Size 32, hàng chính hãng từ Mỹ.",
			categoryID:    14,
			sellerID:      17,
			startingPrice: 500000,
			currentPrice:  650000,
			buyNowPrice:   floatPtr(800000),
			stepPrice:     50000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/levis-511.jpg",
			autoExtend:    true,
			endAt:         now.Add(40 * time.Hour),
			currentBidder: int64Ptr(9),
			bidCount:      3,
			categoryName:  "Quần áo nam",
			parentCatID:   2,
			parentCatName: "Thời trang",
		},
		{
			name:          "Áo sơ mi Oxford nam",
			description:   "Áo sơ mi Oxford cao cấp, form slim fit. Màu xanh navy, size M, L.",
			categoryID:    14,
			sellerID:      18,
			startingPrice: 250000,
			currentPrice:  250000,
			buyNowPrice:   floatPtr(400000),
			stepPrice:     25000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/oxford-shirt.jpg",
			autoExtend:    false,
			endAt:         now.Add(100 * time.Hour),
			currentBidder: nil,
			bidCount:      0,
			categoryName:  "Quần áo nam",
			parentCatID:   2,
			parentCatName: "Thời trang",
		},
		// THỜI TRANG - Quần áo nữ
		{
			name:          "Váy đầm công sở nữ thanh lịch",
			description:   "Váy đầm công sở chất liệu lụa cao cấp, thiết kế thanh lịch. Màu đen, size S, M.",
			categoryID:    15,
			sellerID:      11,
			startingPrice: 400000,
			currentPrice:  400000,
			buyNowPrice:   floatPtr(700000),
			stepPrice:     50000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/office-dress.jpg",
			autoExtend:    false,
			endAt:         now.Add(96 * time.Hour),
			currentBidder: nil,
			bidCount:      0,
			categoryName:  "Quần áo nữ",
			parentCatID:   2,
			parentCatName: "Thời trang",
		},
		// THỜI TRANG - Giày dép
		{
			name:          "Giày Nike Air Force 1 White",
			description:   "Giày Nike Air Force 1 trắng full, hàng chính hãng. Size 40, 41, 42.",
			categoryID:    16,
			sellerID:      10,
			startingPrice: 2000000,
			currentPrice:  2400000,
			buyNowPrice:   floatPtr(3000000),
			stepPrice:     100000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/nike-af1.jpg",
			autoExtend:    true,
			endAt:         now.Add(18 * time.Hour),
			currentBidder: int64Ptr(8),
			bidCount:      4,
			categoryName:  "Giày dép",
			parentCatID:   2,
			parentCatName: "Thời trang",
		},
		{
			name:          "Giày Adidas Ultraboost 22",
			description:   "Adidas Ultraboost 22 màu đen, size 42. Giày chạy bộ cao cấp, mới 99%.",
			categoryID:    16,
			sellerID:      17,
			startingPrice: 1800000,
			currentPrice:  2100000,
			buyNowPrice:   floatPtr(2500000),
			stepPrice:     100000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/ultraboost22.jpg",
			autoExtend:    false,
			endAt:         now.Add(32 * time.Hour),
			currentBidder: int64Ptr(7),
			bidCount:      3,
			categoryName:  "Giày dép",
			parentCatID:   2,
			parentCatName: "Thời trang",
		},
		{
			name:          "Giày Jordan 1 High Bred Toe",
			description:   "Air Jordan 1 High Bred Toe, hàng rep 1:1 cao cấp. Size 41, 42.",
			categoryID:    16,
			sellerID:      18,
			startingPrice: 1200000,
			currentPrice:  1500000,
			buyNowPrice:   floatPtr(1800000),
			stepPrice:     100000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/jordan1-bred.jpg",
			autoExtend:    true,
			endAt:         now.Add(26 * time.Hour),
			currentBidder: int64Ptr(9),
			bidCount:      3,
			categoryName:  "Giày dép",
			parentCatID:   2,
			parentCatName: "Thời trang",
		},
		{
			name:          "Giày Vans Old Skool Black/White",
			description:   "Vans Old Skool classic đen trắng, hàng chính hãng. Size 40, 41.",
			categoryID:    16,
			sellerID:      10,
			startingPrice: 800000,
			currentPrice:  950000,
			buyNowPrice:   floatPtr(1200000),
			stepPrice:     50000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/vans-oldskool.jpg",
			autoExtend:    false,
			endAt:         now.Add(50 * time.Hour),
			currentBidder: int64Ptr(8),
			bidCount:      3,
			categoryName:  "Giày dép",
			parentCatID:   2,
			parentCatName: "Thời trang",
		},
		{
			name:          "Sandal Gucci Slide",
			description:   "Sandal Gucci Slide hàng super fake 1:1. Size 41, 42.",
			categoryID:    16,
			sellerID:      11,
			startingPrice: 600000,
			currentPrice:  750000,
			buyNowPrice:   floatPtr(900000),
			stepPrice:     50000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/gucci-slide.jpg",
			autoExtend:    false,
			endAt:         now.Add(70 * time.Hour),
			currentBidder: int64Ptr(7),
			bidCount:      3,
			categoryName:  "Giày dép",
			parentCatID:   2,
			parentCatName: "Thời trang",
		},
		{
			name:          "Áo blazer nữ Zara",
			description:   "Áo blazer Zara form chuẩn Tây. Màu đen, size S, mới 95%.",
			categoryID:    15,
			sellerID:      18,
			startingPrice: 350000,
			currentPrice:  450000,
			buyNowPrice:   floatPtr(600000),
			stepPrice:     50000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/zara-blazer.jpg",
			autoExtend:    false,
			endAt:         now.Add(34 * time.Hour),
			currentBidder: int64Ptr(8),
			bidCount:      2,
			categoryName:  "Quần áo nữ",
			parentCatID:   2,
			parentCatName: "Thời trang",
		},
		{
			name:          "Set đồ thể thao nữ Adidas",
			description:   "Set áo + quần thể thao Adidas chính hãng. Màu hồng phấn, size M.",
			categoryID:    15,
			sellerID:      11,
			startingPrice: 600000,
			currentPrice:  750000,
			buyNowPrice:   floatPtr(900000),
			stepPrice:     50000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/adidas-set.jpg",
			autoExtend:    false,
			endAt:         now.Add(18 * time.Hour),
			currentBidder: int64Ptr(9),
			bidCount:      3,
			categoryName:  "Quần áo nữ",
			parentCatID:   2,
			parentCatName: "Thời trang",
		},
		{
			name:          "Giày Nike Air Force 1 White",
			description:   "Giày Nike Air Force 1 trắng full, hàng chính hãng. Size 40, 41, 42.",
			categoryID:    16, // Giày dép
			sellerID:      10,
			startingPrice: 2000000,
			currentPrice:  2400000,
			buyNowPrice:   floatPtr(3000000),
			stepPrice:     100000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/nike-af1.jpg",
			autoExtend:    true,
			endAt:         now.Add(18 * time.Hour),
			currentBidder: int64Ptr(8),
			bidCount:      4,
			categoryName:  "Giày dép",
			parentCatID:   2,
			parentCatName: "Thời trang",
		},
		// THỜI TRANG - Túi xách & Đồng hồ
		{
			name:          "Túi xách Michael Kors",
			description:   "Túi xách Michael Kors chính hãng từ Mỹ. Màu đen, da thật, mới 98%.",
			categoryID:    17,
			sellerID:      11,
			startingPrice: 2000000,
			currentPrice:  2400000,
			buyNowPrice:   floatPtr(3000000),
			stepPrice:     100000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/mk-bag.jpg",
			autoExtend:    false,
			endAt:         now.Add(38 * time.Hour),
			currentBidder: int64Ptr(8),
			bidCount:      4,
			categoryName:  "Túi xách",
			parentCatID:   2,
			parentCatName: "Thời trang",
		},
		{
			name:          "Balo Fjallraven Kanken Classic",
			description:   "Balo Fjallraven Kanken màu vàng chanh, hàng chính hãng Thụy Điển. Fullbox.",
			categoryID:    17,
			sellerID:      10,
			startingPrice: 1200000,
			currentPrice:  1400000,
			buyNowPrice:   floatPtr(1800000),
			stepPrice:     100000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/kanken-yellow.jpg",
			autoExtend:    true,
			endAt:         now.Add(56 * time.Hour),
			currentBidder: int64Ptr(9),
			bidCount:      2,
			categoryName:  "Túi xách",
			parentCatID:   2,
			parentCatName: "Thời trang",
		},
		{
			name:          "Đồng hồ Casio G-Shock GA-2100",
			description:   "Casio G-Shock GA-2100 \"CasiOak\" màu đen, hàng chính hãng. Còn bảo hành 11 tháng.",
			categoryID:    18,
			sellerID:      17,
			startingPrice: 2500000,
			currentPrice:  2800000,
			buyNowPrice:   floatPtr(3200000),
			stepPrice:     100000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/gshock-2100.jpg",
			autoExtend:    false,
			endAt:         now.Add(45 * time.Hour),
			currentBidder: int64Ptr(7),
			bidCount:      3,
			categoryName:  "Đồng hồ",
			parentCatID:   2,
			parentCatName: "Thời trang",
		},
		{
			name:          "Đồng hồ Seiko 5 Sports SRPD",
			description:   "Seiko 5 Sports SRPD màu xanh dương, máy cơ automatic. Hàng chính hãng Nhật.",
			categoryID:    18,
			sellerID:      18,
			startingPrice: 4000000,
			currentPrice:  4500000,
			buyNowPrice:   floatPtr(5500000),
			stepPrice:     200000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/seiko5-blue.jpg",
			autoExtend:    true,
			endAt:         now.Add(62 * time.Hour),
			currentBidder: int64Ptr(8),
			bidCount:      2,
			categoryName:  "Đồng hồ",
			parentCatID:   2,
			parentCatName: "Thời trang",
		},

		// ĐỒ GIA DỤNG
		{
			name:          "Bàn làm việc gỗ công nghiệp",
			description:   "Bàn làm việc gỗ công nghiệp cao cấp, kích thước 120x60cm. Màu nâu gỗ tự nhiên.",
			categoryID:    19,
			sellerID:      11,
			startingPrice: 1500000,
			currentPrice:  1800000,
			buyNowPrice:   nil,
			stepPrice:     100000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/work-desk.jpg",
			autoExtend:    false,
			endAt:         now.Add(30 * time.Hour),
			currentBidder: int64Ptr(7),
			bidCount:      3,
			categoryName:  "Nội thất",
			parentCatID:   3,
			parentCatName: "Đồ gia dụng",
		},
		{
			name:          "Ghế gaming E-Dra Hercules",
			description:   "Ghế gaming E-Dra Hercules EGC206, tựa lưng cao, có massage. Màu đen đỏ.",
			categoryID:    19,
			sellerID:      10,
			startingPrice: 2500000,
			currentPrice:  2800000,
			buyNowPrice:   floatPtr(3500000),
			stepPrice:     100000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/gaming-chair.jpg",
			autoExtend:    true,
			endAt:         now.Add(22 * time.Hour),
			currentBidder: int64Ptr(8),
			bidCount:      3,
			categoryName:  "Nội thất",
			parentCatID:   3,
			parentCatName: "Đồ gia dụng",
		},
		{
			name:          "Sofa da 2 chỗ ngồi",
			description:   "Sofa da thật cao cấp màu nâu, 2 chỗ ngồi. Kích thước 160cm, như mới.",
			categoryID:    19,
			sellerID:      18,
			startingPrice: 8000000,
			currentPrice:  9000000,
			buyNowPrice:   floatPtr(12000000),
			stepPrice:     500000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/leather-sofa.jpg",
			autoExtend:    false,
			endAt:         now.Add(80 * time.Hour),
			currentBidder: int64Ptr(9),
			bidCount:      2,
			categoryName:  "Nội thất",
			parentCatID:   3,
			parentCatName: "Đồ gia dụng",
		},
		{
			name:          "Tủ lạnh Samsung Inverter 320L",
			description:   "Tủ lạnh Samsung Inverter 320 lít, tiết kiệm điện. Còn bảo hành 8 tháng.",
			categoryID:    20,
			sellerID:      10,
			startingPrice: 6000000,
			currentPrice:  7000000,
			buyNowPrice:   floatPtr(8500000),
			stepPrice:     200000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/samsung-fridge.jpg",
			autoExtend:    true,
			endAt:         now.Add(5 * time.Minute),
			currentBidder: int64Ptr(8),
			bidCount:      5,
			categoryName:  "Đồ điện gia dụng",
			parentCatID:   3,
			parentCatName: "Đồ gia dụng",
		},
		{
			name:          "Máy giặt LG Inverter 9kg",
			description:   "Máy giặt LG Inverter 9kg, lồng ngang. Tiết kiệm nước và điện, còn bảo hành 1 năm.",
			categoryID:    20,
			sellerID:      11,
			startingPrice: 5000000,
			currentPrice:  5800000,
			buyNowPrice:   floatPtr(7000000),
			stepPrice:     200000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/lg-washer.jpg",
			autoExtend:    false,
			endAt:         now.Add(46 * time.Hour),
			currentBidder: int64Ptr(7),
			bidCount:      4,
			categoryName:  "Đồ điện gia dụng",
			parentCatID:   3,
			parentCatName: "Đồ gia dụng",
		},
		{
			name:          "Lò vi sóng Sharp 25L",
			description:   "Lò vi sóng Sharp 25 lít, công suất 900W. Mới 95%, còn bảo hành 10 tháng.",
			categoryID:    20,
			sellerID:      17,
			startingPrice: 1200000,
			currentPrice:  1400000,
			buyNowPrice:   floatPtr(1800000),
			stepPrice:     100000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/sharp-microwave.jpg",
			autoExtend:    true,
			endAt:         now.Add(35 * time.Hour),
			currentBidder: int64Ptr(9),
			bidCount:      2,
			categoryName:  "Đồ điện gia dụng",
			parentCatID:   3,
			parentCatName: "Đồ gia dụng",
		},
		{
			name:          "Nồi cơm điện Zojirushi 1.8L",
			description:   "Nồi cơm điện Zojirushi Nhật Bản 1.8L, công nghệ IH. Hàng xách tay Nhật, mới 100%.",
			categoryID:    21,
			sellerID:      10,
			startingPrice: 8000000,
			currentPrice:  8500000,
			buyNowPrice:   floatPtr(10000000),
			stepPrice:     200000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/zojirushi.jpg",
			autoExtend:    false,
			endAt:         now.Add(68 * time.Hour),
			currentBidder: int64Ptr(8),
			bidCount:      2,
			categoryName:  "Đồ dùng nhà bếp",
			parentCatID:   3,
			parentCatName: "Đồ gia dụng",
		},
		{
			name:          "Set dao bếp Nhật Bản",
			description:   "Set 5 dao bếp Nhật Bản, thép không gỉ cao cấp. Kèm khay gỗ đựng dao.",
			categoryID:    21,
			sellerID:      18,
			startingPrice: 1500000,
			currentPrice:  1700000,
			buyNowPrice:   floatPtr(2200000),
			stepPrice:     100000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/knife-set.jpg",
			autoExtend:    true,
			endAt:         now.Add(41 * time.Hour),
			currentBidder: int64Ptr(7),
			bidCount:      2,
			categoryName:  "Đồ dùng nhà bếp",
			parentCatID:   3,
			parentCatName: "Đồ gia dụng",
		},
		// XE CỘ & SÁCH
		{
			name:          "Xe máy Honda Wave RSX 2022",
			description:   "Honda Wave RSX 2022, màu đen đỏ. Xe đẹp, máy êm, chạy 15000km.",
			categoryID:    23,
			sellerID:      17,
			startingPrice: 18000000,
			currentPrice:  20000000,
			buyNowPrice:   floatPtr(22000000),
			stepPrice:     500000,
			status:        "FINISHED",
			thumbnailURL:  "https://example.com/honda-wave.jpg",
			autoExtend:    false,
			endAt:         now.Add(-24 * time.Hour),
			currentBidder: int64Ptr(7),
			bidCount:      4,
			categoryName:  "Xe máy",
			parentCatID:   4,
			parentCatName: "Xe cộ",
		},
		{
			name:          "Yamaha Exciter 155 VVA 2023",
			description:   "Exciter 155 VVA màu xanh GP, đời 2023. Xe zin 100%, chạy 8000km, còn bảo hành.",
			categoryID:    23,
			sellerID:      10,
			startingPrice: 38000000,
			currentPrice:  41000000,
			buyNowPrice:   floatPtr(45000000),
			stepPrice:     500000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/exciter155.jpg",
			autoExtend:    true,
			endAt:         now.Add(15 * time.Hour),
			currentBidder: int64Ptr(8),
			bidCount:      6,
			categoryName:  "Xe máy",
			parentCatID:   4,
			parentCatName: "Xe cộ",
		},
		// THỂ THAO & ĐỒ CHƠI
		{
			name:          "Bóng đá FIFA Quality Pro",
			description:   "Bóng đá FIFA Quality Pro, size 5 chuẩn thi đấu. Hàng chính hãng.",
			categoryID:    28,
			sellerID:      10,
			startingPrice: 400000,
			currentPrice:  550000,
			buyNowPrice:   floatPtr(700000),
			stepPrice:     50000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/fifa-ball.jpg",
			autoExtend:    false,
			endAt:         now.Add(8 * time.Hour),
			currentBidder: int64Ptr(9),
			bidCount:      3,
			categoryName:  "Dụng cụ thể thao",
			parentCatID:   6,
			parentCatName: "Thể thao & Du lịch",
		},
		{
			name:          "Vợt cầu lông Yonex Astrox 99 Pro",
			description:   "Vợt cầu lông Yonex Astrox 99 Pro chính hãng, đã căng cước Yonex BG80.",
			categoryID:    28,
			sellerID:      11,
			startingPrice: 3500000,
			currentPrice:  3800000,
			buyNowPrice:   floatPtr(4500000),
			stepPrice:     100000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/yonex-99pro.jpg",
			autoExtend:    true,
			endAt:         now.Add(50 * time.Hour),
			currentBidder: int64Ptr(7),
			bidCount:      3,
			categoryName:  "Dụng cụ thể thao",
			parentCatID:   6,
			parentCatName: "Thể thao & Du lịch",
		},
		{
			name:          "Xe đạp Road Giant TCR",
			description:   "Xe đạp đua Giant TCR Advanced, khung carbon. Group Shimano 105, mới 95%.",
			categoryID:    28,
			sellerID:      18,
			startingPrice: 18000000,
			currentPrice:  20000000,
			buyNowPrice:   floatPtr(25000000),
			stepPrice:     500000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/giant-tcr.jpg",
			autoExtend:    false,
			endAt:         now.Add(88 * time.Hour),
			currentBidder: int64Ptr(8),
			bidCount:      4,
			categoryName:  "Dụng cụ thể thao",
			parentCatID:   6,
			parentCatName: "Thể thao & Du lịch",
		},
		{
			name:          "Lego Technic Lamborghini Sián",
			description:   "Bộ Lego Technic Lamborghini Sián FKP 37, fullbox, chưa mở seal.",
			categoryID:    30,
			sellerID:      11,
			startingPrice: 8000000,
			currentPrice:  8000000,
			buyNowPrice:   floatPtr(12000000),
			stepPrice:     200000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/lego-lambo.jpg",
			autoExtend:    true,
			endAt:         now.Add(120 * time.Hour),
			currentBidder: nil,
			bidCount:      0,
			categoryName:  "Đồ chơi trẻ em",
			parentCatID:   7,
			parentCatName: "Đồ chơi & Trẻ em",
		},
		{
			name:          "Mô hình Gundam RG Nu Gundam",
			description:   "Mô hình Gundam RG Nu Gundam tỷ lệ 1/144, hàng Bandai chính hãng. Chưa lắp ráp.",
			categoryID:    30,
			sellerID:      10,
			startingPrice: 600000,
			currentPrice:  700000,
			buyNowPrice:   floatPtr(900000),
			stepPrice:     50000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/gundam-rg-nu.jpg",
			autoExtend:    false,
			endAt:         now.Add(65 * time.Hour),
			currentBidder: int64Ptr(9),
			bidCount:      2,
			categoryName:  "Đồ chơi trẻ em",
			parentCatID:   7,
			parentCatName: "Đồ chơi & Trẻ em",
		},
		{
			name:          "Xe điều khiển Traxxas X-Maxx",
			description:   "Xe điều khiển Monster Truck Traxxas X-Maxx 1/5, brushless. Hàng nhập Mỹ.",
			categoryID:    30,
			sellerID:      17,
			startingPrice: 12000000,
			currentPrice:  13000000,
			buyNowPrice:   floatPtr(16000000),
			stepPrice:     500000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/traxxas-xmaxx.jpg",
			autoExtend:    true,
			endAt:         now.Add(72 * time.Hour),
			currentBidder: int64Ptr(7),
			bidCount:      2,
			categoryName:  "Đồ chơi trẻ em",
			parentCatID:   7,
			parentCatName: "Đồ chơi & Trẻ em",
		},
		// MỸ PHẨM & LÀM ĐẸP
		{
			name:          "Kem chống nắng La Roche-Posay SPF50+",
			description:   "Kem chống nắng La Roche-Posay Anthelios SPF50+ 50ml. Hàng Pháp chính hãng.",
			categoryID:    33,
			sellerID:      10,
			startingPrice: 300000,
			currentPrice:  380000,
			buyNowPrice:   floatPtr(450000),
			stepPrice:     20000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/laroche-sunscreen.jpg",
			autoExtend:    false,
			endAt:         now.Add(2 * time.Minute),
			currentBidder: int64Ptr(7),
			bidCount:      4,
			categoryName:  "Chăm sóc da",
			parentCatID:   8,
			parentCatName: "Mỹ phẩm & Làm đẹp",
		},
		{
			name:          "Serum Vitamin C The Ordinary",
			description:   "The Ordinary Vitamin C Suspension 23% + HA Spheres 2% 30ml. Hàng Canada.",
			categoryID:    33,
			sellerID:      11,
			startingPrice: 200000,
			currentPrice:  250000,
			buyNowPrice:   floatPtr(350000),
			stepPrice:     20000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/ordinary-vitc.jpg",
			autoExtend:    false,
			endAt:         now.Add(33 * time.Hour),
			currentBidder: int64Ptr(8),
			bidCount:      2,
			categoryName:  "Chăm sóc da",
			parentCatID:   8,
			parentCatName: "Mỹ phẩm & Làm đẹp",
		},
		{
			name:          "Son Dior Addict Lip Glow",
			description:   "Son dưỡng Dior Addict Lip Glow màu 001 Pink. Hàng Pháp chính hãng, mới 95%.",
			categoryID:    32,
			sellerID:      18,
			startingPrice: 600000,
			currentPrice:  700000,
			buyNowPrice:   floatPtr(850000),
			stepPrice:     50000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/dior-lipglow.jpg",
			autoExtend:    true,
			endAt:         now.Add(27 * time.Hour),
			currentBidder: int64Ptr(9),
			bidCount:      2,
			categoryName:  "Mỹ phẩm",
			parentCatID:   8,
			parentCatName: "Mỹ phẩm & Làm đẹp",
		},
		{
			name:          "Set Innisfree Green Tea Seed Serum",
			description:   "Set dưỡng da Innisfree Green Tea Seed Serum gồm 4 món. Hàng Hàn Quốc chính hãng.",
			categoryID:    33,
			sellerID:      10,
			startingPrice: 800000,
			currentPrice:  900000,
			buyNowPrice:   floatPtr(1100000),
			stepPrice:     50000,
			status:        "ACTIVE",
			thumbnailURL:  "https://example.com/innisfree-set.jpg",
			autoExtend:    false,
			endAt:         now.Add(55 * time.Hour),
			currentBidder: int64Ptr(7),
			bidCount:      2,
			categoryName:  "Chăm sóc da",
			parentCatID:   8,
			parentCatName: "Mỹ phẩm & Làm đẹp",
		},
	}

	for i, p := range products {
		// Determine order_created based on status
		orderCreated := p.status == "FINISHED"
		
		var productID int64
		err := db.QueryRowContext(ctx, `
			INSERT INTO products (
				name, description, category_id, seller_id, starting_price, current_price, 
				buy_now_price, step_price, status, thumbnail_url, auto_extend, end_at, 
				created_at, current_bidder, bid_count, category_name, parent_category_id, parent_category_name,
				order_created, sent_email
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
			RETURNING id
		`, p.name, p.description, p.categoryID, p.sellerID, p.startingPrice, p.currentPrice,
			p.buyNowPrice, p.stepPrice, p.status, p.thumbnailURL, p.autoExtend, p.endAt,
			now, p.currentBidder, p.bidCount, p.categoryName, p.parentCatID, p.parentCatName,
			orderCreated, false).Scan(&productID)
		if err != nil {
			return nil, fmt.Errorf("error inserting product %s: %v", p.name, err)
		}
		
		// Save the mapping: index (1-based) -> actual product ID
		productIDMap[i+1] = productID
	}

	log.Println("Products seeded successfully")
	return productIDMap, nil
}

func seedProductImages(ctx context.Context, db *sql.DB, productIDMap map[int]int64) error {
	log.Println("Seeding product images...")

	images := []struct {
		productIndex int    // index in products array (1-based)
		imageURL     string
	}{
		// Product 1-4: Điện thoại
		{1, "https://example.com/iphone15promax-1.jpg"},
		{1, "https://example.com/iphone15promax-2.jpg"},
		{1, "https://example.com/iphone15promax-3.jpg"},
		{2, "https://example.com/s24ultra-1.jpg"},
		{2, "https://example.com/s24ultra-2.jpg"},
		{2, "https://example.com/s24ultra-3.jpg"},
		{3, "https://example.com/iphone14pro-1.jpg"},
		{3, "https://example.com/iphone14pro-2.jpg"},
		{4, "https://example.com/xiaomi14ultra-1.jpg"},
		{4, "https://example.com/xiaomi14ultra-2.jpg"},

		// Product 5-9: Laptop
		{5, "https://example.com/macbook-m3-1.jpg"},
		{5, "https://example.com/macbook-m3-2.jpg"},
		{5, "https://example.com/macbook-m3-3.jpg"},
		{6, "https://example.com/dell-xps15-1.jpg"},
		{6, "https://example.com/dell-xps15-2.jpg"},
		{7, "https://example.com/asus-rog-1.jpg"},
		{7, "https://example.com/asus-rog-2.jpg"},
		{8, "https://example.com/thinkpad-1.jpg"},
		{9, "https://example.com/macair-m2-1.jpg"},
		{9, "https://example.com/macair-m2-2.jpg"},

		// Product 10-12: Tablet
		{10, "https://example.com/ipad-pro-1.jpg"},
		{10, "https://example.com/ipad-pro-2.jpg"},
		{11, "https://example.com/ipad-air5-1.jpg"},
		{12, "https://example.com/tab-s9-1.jpg"},
		{12, "https://example.com/tab-s9-2.jpg"},

		// Product 13-16: Tai nghe
		{13, "https://example.com/airpods-1.jpg"},
		{14, "https://example.com/sony-1000xm5-1.jpg"},
		{14, "https://example.com/sony-1000xm5-2.jpg"},
		{15, "https://example.com/bose-qc45-1.jpg"},
		{16, "https://example.com/airpods-max-1.jpg"},
		{16, "https://example.com/airpods-max-2.jpg"},

		// Product 17-19: Smart Watch
		{17, "https://example.com/watch-1.jpg"},
		{17, "https://example.com/watch-2.jpg"},
		{18, "https://example.com/galaxy-watch-1.jpg"},
		{19, "https://example.com/garmin-1.jpg"},

		// Product 20-23: Quần áo nam
		{20, "https://example.com/bomber-1.jpg"},
		{21, "https://example.com/polo-1.jpg"},
		{22, "https://example.com/levis-1.jpg"},
		{23, "https://example.com/oxford-1.jpg"},

		// Product 24-27: Quần áo nữ
		{24, "https://example.com/dress-1.jpg"},
		{25, "https://example.com/coat-1.jpg"},
		{25, "https://example.com/coat-2.jpg"},
		{26, "https://example.com/blazer-1.jpg"},
		{27, "https://example.com/adidas-set-1.jpg"},

		// Product 28-32: Giày dép
		{28, "https://example.com/nike-1.jpg"},
		{28, "https://example.com/nike-2.jpg"},
		{29, "https://example.com/ultraboost-1.jpg"},
		{30, "https://example.com/jordan-1.jpg"},
		{31, "https://example.com/vans-1.jpg"},
		{32, "https://example.com/gucci-slide-1.jpg"},

		// Product 33-36: Túi xách & Đồng hồ
		{33, "https://example.com/mk-bag-1.jpg"},
		{33, "https://example.com/mk-bag-2.jpg"},
		{34, "https://example.com/kanken-1.jpg"},
		{35, "https://example.com/gshock-1.jpg"},
		{36, "https://example.com/seiko5-1.jpg"},

		// Product 37-40: Đồ gia dụng
		{37, "https://example.com/desk-1.jpg"},
		{38, "https://example.com/gaming-chair-1.jpg"},
		{38, "https://example.com/gaming-chair-2.jpg"},
		{39, "https://example.com/sofa-1.jpg"},
		{39, "https://example.com/sofa-2.jpg"},
		{40, "https://example.com/fridge-1.jpg"},
		{41, "https://example.com/lg-washer-1.jpg"},
		{42, "https://example.com/microwave-1.jpg"},
		{43, "https://example.com/zojirushi-1.jpg"},
		{44, "https://example.com/knife-set-1.jpg"},

		// Product 45-47: Xe máy
		{45, "https://example.com/wave-1.jpg"},
		{45, "https://example.com/wave-2.jpg"},
		{46, "https://example.com/exciter-1.jpg"},
		{46, "https://example.com/exciter-2.jpg"},
		{47, "https://example.com/sh-mode-1.jpg"},

		// Product 48-50: Sách
		{48, "https://example.com/book-1.jpg"},
		{49, "https://example.com/harry-potter-1.jpg"},
		{50, "https://example.com/nha-gia-kim-1.jpg"},

		// Product 51-53: Thể thao
		{51, "https://example.com/ball-1.jpg"},
		{52, "https://example.com/yonex-1.jpg"},
		{53, "https://example.com/giant-tcr-1.jpg"},
		{53, "https://example.com/giant-tcr-2.jpg"},

		// Product 54-56: Đồ chơi
		{54, "https://example.com/lego-1.jpg"},
		{54, "https://example.com/lego-2.jpg"},
		{55, "https://example.com/gundam-1.jpg"},
		{56, "https://example.com/traxxas-1.jpg"},
		{56, "https://example.com/traxxas-2.jpg"},
	}

	for _, img := range images {
		// Get actual product ID from map
		productID, exists := productIDMap[img.productIndex]
		if !exists {
			return fmt.Errorf("product index %d not found in product ID map", img.productIndex)
		}
		
		_, err := db.ExecContext(ctx, `
			INSERT INTO product_images (product_id, image_url)
			VALUES ($1, $2)
		`, productID, img.imageURL)
		if err != nil {
			return fmt.Errorf("error inserting product image for product index %d: %v", img.productIndex, err)
		}
	}

	log.Println("Product images seeded successfully")
	return nil
}

func seedWatchList(ctx context.Context, db *sql.DB) error {
	log.Println("Seeding watch list...")

	now := time.Now()

	watchList := []struct {
		userID    int64
		productID int64
	}{
		// User 7 (Alice - Bidder)
		{7, 1}, {7, 2}, {7, 5}, {7, 10}, {7, 13}, {7, 28}, {7, 40}, {7, 46},
		{7, 17}, {7, 33}, {7, 52}, {7, 54},

		// User 8 (Bob - Bidder)
		{8, 1}, {8, 3}, {8, 6}, {8, 8}, {8, 14}, {8, 20}, {8, 29}, {8, 41},
		{8, 45}, {8, 51}, {8, 56},

		// User 9 (Charlie - Bidder)
		{9, 2}, {9, 4}, {9, 7}, {9, 11}, {9, 15}, {9, 24}, {9, 30}, {9, 39},
		{9, 47}, {9, 53}, {9, 57},

		// User 21 (Trí - Bidder)
		{21, 1}, {21, 5}, {21, 10}, {21, 28}, {21, 40}, {21, 46}, {21, 54},

		// User 23 (Test User - Bidder)
		{23, 2}, {23, 6}, {23, 13}, {23, 20}, {23, 51},

		// User 24 (Tri Ngo - Bidder)
		{24, 1}, {24, 10}, {24, 17}, {24, 28}, {24, 40},
	}

	for _, w := range watchList {
		_, err := db.ExecContext(ctx, `
			INSERT INTO watch_list (user_id, product_id, created_at)
			VALUES ($1, $2, $3)
		`, w.userID, w.productID, now)
		if err != nil {
			return fmt.Errorf("error inserting watch list: %v", err)
		}
	}

	log.Println("Watch list seeded successfully")
	return nil
}

func seedBiddingHistory(ctx context.Context, db *sql.DB) error {
	log.Println("Seeding bidding history...")

	now := time.Now()

	biddingHistory := []struct {
		amount    float64
		bidderID  int64
		productID int64
		status    string
		reason    *string
		requestID string
	}{
		// Product 1: iPhone 15 Pro Max (5 bids)
		{25000000, 7, 1, "SUCCESS", nil, "BID-1-1"},
		{25500000, 8, 1, "SUCCESS", nil, "BID-1-2"},
		{26000000, 7, 1, "SUCCESS", nil, "BID-1-3"},
		{26500000, 9, 1, "SUCCESS", nil, "BID-1-4"},
		{27500000, 7, 1, "SUCCESS", nil, "BID-1-5"},

		// Product 2: Samsung S24 Ultra (4 bids)
		{24000000, 8, 2, "SUCCESS", nil, "BID-2-1"},
		{24500000, 7, 2, "SUCCESS", nil, "BID-2-2"},
		{25000000, 9, 2, "SUCCESS", nil, "BID-2-3"},
		{26000000, 8, 2, "SUCCESS", nil, "BID-2-4"},

		// Product 3: iPhone 14 Pro (4 bids)
		{18000000, 7, 3, "SUCCESS", nil, "BID-3-1"},
		{18500000, 8, 3, "SUCCESS", nil, "BID-3-2"},
		{19000000, 9, 3, "SUCCESS", nil, "BID-3-3"},
		{20000000, 9, 3, "SUCCESS", nil, "BID-3-4"},

		// Product 5: MacBook Pro M3 (6 bids)
		{35000000, 7, 5, "SUCCESS", nil, "BID-5-1"},
		{35500000, 8, 5, "SUCCESS", nil, "BID-5-2"},
		{36000000, 7, 5, "SUCCESS", nil, "BID-5-3"},
		{36500000, 9, 5, "SUCCESS", nil, "BID-5-4"},
		{37000000, 8, 5, "SUCCESS", nil, "BID-5-5"},
		{38000000, 8, 5, "SUCCESS", nil, "BID-5-6"},

		// Product 6: Dell XPS 15 (6 bids)
		{30000000, 7, 6, "SUCCESS", nil, "BID-6-1"},
		{30500000, 8, 6, "SUCCESS", nil, "BID-6-2"},
		{31000000, 9, 6, "SUCCESS", nil, "BID-6-3"},
		{31500000, 7, 6, "SUCCESS", nil, "BID-6-4"},
		{32000000, 8, 6, "SUCCESS", nil, "BID-6-5"},
		{33000000, 7, 6, "SUCCESS", nil, "BID-6-6"},

		// Product 7: ASUS ROG (4 bids)
		{28000000, 9, 7, "SUCCESS", nil, "BID-7-1"},
		{28500000, 7, 7, "SUCCESS", nil, "BID-7-2"},
		{29000000, 8, 7, "SUCCESS", nil, "BID-7-3"},
		{30000000, 8, 7, "SUCCESS", nil, "BID-7-4"},

		// Product 8: Lenovo ThinkPad (4 bids)
		{25000000, 7, 8, "SUCCESS", nil, "BID-8-1"},
		{25500000, 8, 8, "SUCCESS", nil, "BID-8-2"},
		{26000000, 9, 8, "SUCCESS", nil, "BID-8-3"},
		{27000000, 9, 8, "SUCCESS", nil, "BID-8-4"},

		// Product 9: MacBook Air M2 (4 bids)
		{20000000, 7, 9, "SUCCESS", nil, "BID-9-1"},
		{20500000, 8, 9, "SUCCESS", nil, "BID-9-2"},
		{21000000, 9, 9, "SUCCESS", nil, "BID-9-3"},
		{22000000, 7, 9, "SUCCESS", nil, "BID-9-4"},

		// Product 10: iPad Pro M2 (10 bids)
		{18000000, 7, 10, "SUCCESS", nil, "BID-10-1"},
		{18200000, 8, 10, "SUCCESS", nil, "BID-10-2"},
		{18400000, 7, 10, "SUCCESS", nil, "BID-10-3"},
		{18600000, 9, 10, "SUCCESS", nil, "BID-10-4"},
		{18800000, 7, 10, "SUCCESS", nil, "BID-10-5"},
		{19000000, 8, 10, "SUCCESS", nil, "BID-10-6"},
		{19200000, 7, 10, "SUCCESS", nil, "BID-10-7"},
		{19400000, 9, 10, "SUCCESS", nil, "BID-10-8"},
		{19600000, 8, 10, "SUCCESS", nil, "BID-10-9"},
		{20000000, 7, 10, "SUCCESS", nil, "BID-10-10"},

		// Product 11: iPad Air 5 (7 bids)
		{12000000, 7, 11, "SUCCESS", nil, "BID-11-1"},
		{12200000, 8, 11, "SUCCESS", nil, "BID-11-2"},
		{12400000, 9, 11, "SUCCESS", nil, "BID-11-3"},
		{12600000, 7, 11, "SUCCESS", nil, "BID-11-4"},
		{12800000, 8, 11, "SUCCESS", nil, "BID-11-5"},
		{13000000, 9, 11, "SUCCESS", nil, "BID-11-6"},
		{13500000, 8, 11, "SUCCESS", nil, "BID-11-7"},

		// Product 12: Tab S9+ (5 bids)
		{15000000, 7, 12, "SUCCESS", nil, "BID-12-1"},
		{15200000, 8, 12, "SUCCESS", nil, "BID-12-2"},
		{15400000, 9, 12, "SUCCESS", nil, "BID-12-3"},
		{15600000, 7, 12, "SUCCESS", nil, "BID-12-4"},
		{16000000, 9, 12, "SUCCESS", nil, "BID-12-5"},

		// Product 13: AirPods Pro 2 (5 bids)
		{5000000, 7, 13, "SUCCESS", nil, "BID-13-1"},
		{5100000, 8, 13, "SUCCESS", nil, "BID-13-2"},
		{5200000, 9, 13, "SUCCESS", nil, "BID-13-3"},
		{5300000, 8, 13, "SUCCESS", nil, "BID-13-4"},
		{5500000, 9, 13, "SUCCESS", nil, "BID-13-5"},

		// Product 14: Sony WH-1000XM5 (4 bids)
		{6000000, 8, 14, "SUCCESS", nil, "BID-14-1"},
		{6200000, 7, 14, "SUCCESS", nil, "BID-14-2"},
		{6400000, 9, 14, "SUCCESS", nil, "BID-14-3"},
		{6800000, 7, 14, "SUCCESS", nil, "BID-14-4"},

		// Product 15: Bose QC45 (5 bids)
		{4000000, 7, 15, "SUCCESS", nil, "BID-15-1"},
		{4100000, 9, 15, "SUCCESS", nil, "BID-15-2"},
		{4200000, 8, 15, "SUCCESS", nil, "BID-15-3"},
		{4300000, 7, 15, "SUCCESS", nil, "BID-15-4"},
		{4500000, 8, 15, "SUCCESS", nil, "BID-15-5"},

		// Product 16: AirPods Max (5 bids)
		{10000000, 7, 16, "SUCCESS", nil, "BID-16-1"},
		{10200000, 8, 16, "SUCCESS", nil, "BID-16-2"},
		{10400000, 9, 16, "SUCCESS", nil, "BID-16-3"},
		{10600000, 7, 16, "SUCCESS", nil, "BID-16-4"},
		{11000000, 9, 16, "SUCCESS", nil, "BID-16-5"},

		// Product 17: Apple Watch S9 (5 bids)
		{8000000, 8, 17, "SUCCESS", nil, "BID-17-1"},
		{8200000, 7, 17, "SUCCESS", nil, "BID-17-2"},
		{8400000, 9, 17, "SUCCESS", nil, "BID-17-3"},
		{8600000, 7, 17, "SUCCESS", nil, "BID-17-4"},
		{9000000, 8, 17, "SUCCESS", nil, "BID-17-5"},

		// Product 18: Galaxy Watch 6 (2 bids)
		{6000000, 9, 18, "SUCCESS", nil, "BID-18-1"},
		{6500000, 7, 18, "SUCCESS", nil, "BID-18-2"},

		// Product 19: Garmin Fenix 7X (2 bids)
		{15000000, 7, 19, "SUCCESS", nil, "BID-19-1"},
		{16000000, 8, 19, "SUCCESS", nil, "BID-19-2"},

		// Product 20: Áo Bomber (3 bids)
		{300000, 7, 20, "SUCCESS", nil, "BID-20-1"},
		{350000, 8, 20, "SUCCESS", nil, "BID-20-2"},
		{450000, 7, 20, "SUCCESS", nil, "BID-20-3"},

		// Product 21: Polo Ralph Lauren (3 bids)
		{400000, 9, 21, "SUCCESS", nil, "BID-21-1"},
		{450000, 7, 21, "SUCCESS", nil, "BID-21-2"},
		{550000, 8, 21, "SUCCESS", nil, "BID-21-3"},

		// Product 22: Levi's 511 (3 bids)
		{500000, 8, 22, "SUCCESS", nil, "BID-22-1"},
		{550000, 7, 22, "SUCCESS", nil, "BID-22-2"},
		{650000, 9, 22, "SUCCESS", nil, "BID-22-3"},

		// Product 25: Áo khoác dạ (3 bids)
		{800000, 8, 25, "SUCCESS", nil, "BID-25-1"},
		{850000, 9, 25, "SUCCESS", nil, "BID-25-2"},
		{950000, 7, 25, "SUCCESS", nil, "BID-25-3"},

		// Product 26: Blazer Zara (2 bids)
		{350000, 7, 26, "SUCCESS", nil, "BID-26-1"},
		{450000, 8, 26, "SUCCESS", nil, "BID-26-2"},

		// Product 27: Set Adidas (3 bids)
		{600000, 8, 27, "SUCCESS", nil, "BID-27-1"},
		{650000, 7, 27, "SUCCESS", nil, "BID-27-2"},
		{750000, 9, 27, "SUCCESS", nil, "BID-27-3"},

		// Product 28: Nike AF1 (4 bids)
		{2000000, 8, 28, "SUCCESS", nil, "BID-28-1"},
		{2100000, 7, 28, "SUCCESS", nil, "BID-28-2"},
		{2200000, 9, 28, "SUCCESS", nil, "BID-28-3"},
		{2400000, 8, 28, "SUCCESS", nil, "BID-28-4"},

		// Product 29: Ultraboost 22 (3 bids)
		{1800000, 9, 29, "SUCCESS", nil, "BID-29-1"},
		{1900000, 8, 29, "SUCCESS", nil, "BID-29-2"},
		{2100000, 7, 29, "SUCCESS", nil, "BID-29-3"},

		// Product 30: Jordan 1 (3 bids)
		{1200000, 7, 30, "SUCCESS", nil, "BID-30-1"},
		{1300000, 8, 30, "SUCCESS", nil, "BID-30-2"},
		{1500000, 9, 30, "SUCCESS", nil, "BID-30-3"},

		// Product 31: Vans Old Skool (3 bids)
		{800000, 9, 31, "SUCCESS", nil, "BID-31-1"},
		{850000, 7, 31, "SUCCESS", nil, "BID-31-2"},
		{950000, 8, 31, "SUCCESS", nil, "BID-31-3"},

		// Product 32: Gucci Slide (3 bids)
		{600000, 8, 32, "SUCCESS", nil, "BID-32-1"},
		{650000, 9, 32, "SUCCESS", nil, "BID-32-2"},
		{750000, 7, 32, "SUCCESS", nil, "BID-32-3"},

		// Product 33: Túi MK (4 bids)
		{2000000, 7, 33, "SUCCESS", nil, "BID-33-1"},
		{2100000, 9, 33, "SUCCESS", nil, "BID-33-2"},
		{2200000, 7, 33, "SUCCESS", nil, "BID-33-3"},
		{2400000, 8, 33, "SUCCESS", nil, "BID-33-4"},

		// Product 34: Kanken (2 bids)
		{1200000, 8, 34, "SUCCESS", nil, "BID-34-1"},
		{1400000, 9, 34, "SUCCESS", nil, "BID-34-2"},

		// Product 35: G-Shock (3 bids)
		{2500000, 8, 35, "SUCCESS", nil, "BID-35-1"},
		{2600000, 9, 35, "SUCCESS", nil, "BID-35-2"},
		{2800000, 7, 35, "SUCCESS", nil, "BID-35-3"},

		// Product 36: Seiko 5 (2 bids)
		{4000000, 7, 36, "SUCCESS", nil, "BID-36-1"},
		{4500000, 8, 36, "SUCCESS", nil, "BID-36-2"},

		// Product 37: Bàn làm việc (3 bids)
		{1500000, 7, 37, "SUCCESS", nil, "BID-37-1"},
		{1600000, 8, 37, "SUCCESS", nil, "BID-37-2"},
		{1800000, 7, 37, "SUCCESS", nil, "BID-37-3"},

		// Product 38: Ghế gaming (3 bids)
		{2500000, 9, 38, "SUCCESS", nil, "BID-38-1"},
		{2600000, 7, 38, "SUCCESS", nil, "BID-38-2"},
		{2800000, 8, 38, "SUCCESS", nil, "BID-38-3"},

		// Product 39: Sofa (2 bids)
		{8000000, 8, 39, "SUCCESS", nil, "BID-39-1"},
		{9000000, 9, 39, "SUCCESS", nil, "BID-39-2"},

		// Product 40: Tủ lạnh (5 bids)
		{6000000, 7, 40, "SUCCESS", nil, "BID-40-1"},
		{6200000, 8, 40, "SUCCESS", nil, "BID-40-2"},
		{6400000, 9, 40, "SUCCESS", nil, "BID-40-3"},
		{6600000, 8, 40, "SUCCESS", nil, "BID-40-4"},
		{7000000, 8, 40, "SUCCESS", nil, "BID-40-5"},

		// Product 41: Máy giặt LG (4 bids)
		{5000000, 8, 41, "SUCCESS", nil, "BID-41-1"},
		{5200000, 9, 41, "SUCCESS", nil, "BID-41-2"},
		{5400000, 8, 41, "SUCCESS", nil, "BID-41-3"},
		{5800000, 7, 41, "SUCCESS", nil, "BID-41-4"},

		// Product 42: Lò vi sóng (2 bids)
		{1200000, 8, 42, "SUCCESS", nil, "BID-42-1"},
		{1400000, 9, 42, "SUCCESS", nil, "BID-42-2"},

		// Product 43: Nồi cơm Zojirushi (2 bids)
		{8000000, 7, 43, "SUCCESS", nil, "BID-43-1"},
		{8500000, 8, 43, "SUCCESS", nil, "BID-43-2"},

		// Product 44: Set dao (2 bids)
		{1500000, 9, 44, "SUCCESS", nil, "BID-44-1"},
		{1700000, 7, 44, "SUCCESS", nil, "BID-44-2"},

		// Product 45: Honda Wave (FINISHED - 4 bids)
		{18000000, 7, 45, "SUCCESS", nil, "BID-45-1"},
		{18500000, 8, 45, "SUCCESS", nil, "BID-45-2"},
		{19000000, 7, 45, "SUCCESS", nil, "BID-45-3"},
		{20000000, 7, 45, "SUCCESS", nil, "BID-45-4"},

		// Product 46: Exciter 155 (6 bids)
		{38000000, 7, 46, "SUCCESS", nil, "BID-46-1"},
		{38500000, 9, 46, "SUCCESS", nil, "BID-46-2"},
		{39000000, 7, 46, "SUCCESS", nil, "BID-46-3"},
		{39500000, 8, 46, "SUCCESS", nil, "BID-46-4"},
		{40000000, 9, 46, "SUCCESS", nil, "BID-46-5"},
		{41000000, 8, 46, "SUCCESS", nil, "BID-46-6"},

		// Product 47: SH Mode (4 bids)
		{48000000, 8, 47, "SUCCESS", nil, "BID-47-1"},
		{48500000, 7, 47, "SUCCESS", nil, "BID-47-2"},
		{49000000, 9, 47, "SUCCESS", nil, "BID-47-3"},
		{50000000, 9, 47, "SUCCESS", nil, "BID-47-4"},

		// Product 48: Đắc Nhân Tâm (FINISHED - 3 bids)
		{50000, 8, 48, "SUCCESS", nil, "BID-48-1"},
		{60000, 7, 48, "SUCCESS", nil, "BID-48-2"},
		{80000, 8, 48, "SUCCESS", nil, "BID-48-3"},

		// Product 49: Harry Potter (3 bids)
		{800000, 8, 49, "SUCCESS", nil, "BID-49-1"},
		{850000, 9, 49, "SUCCESS", nil, "BID-49-2"},
		{950000, 7, 49, "SUCCESS", nil, "BID-49-3"},

		// Product 51: Bóng đá (3 bids)
		{400000, 9, 51, "SUCCESS", nil, "BID-51-1"},
		{450000, 8, 51, "SUCCESS", nil, "BID-51-2"},
		{550000, 9, 51, "SUCCESS", nil, "BID-51-3"},

		// Product 52: Vợt Yonex (3 bids)
		{3500000, 8, 52, "SUCCESS", nil, "BID-52-1"},
		{3600000, 9, 52, "SUCCESS", nil, "BID-52-2"},
		{3800000, 7, 52, "SUCCESS", nil, "BID-52-3"},

		// Product 53: Xe đạp Giant (4 bids)
		{18000000, 7, 53, "SUCCESS", nil, "BID-53-1"},
		{18500000, 9, 53, "SUCCESS", nil, "BID-53-2"},
		{19000000, 7, 53, "SUCCESS", nil, "BID-53-3"},
		{20000000, 8, 53, "SUCCESS", nil, "BID-53-4"},

		// Product 55: Gundam (2 bids)
		{600000, 8, 55, "SUCCESS", nil, "BID-55-1"},
		{700000, 9, 55, "SUCCESS", nil, "BID-55-2"},

		// Product 56: Traxxas X-Maxx (2 bids)
		{12000000, 9, 56, "SUCCESS", nil, "BID-56-1"},
		{13000000, 7, 56, "SUCCESS", nil, "BID-56-2"},

		// Product 57: Kem chống nắng (4 bids)
		{300000, 7, 57, "SUCCESS", nil, "BID-57-1"},
		{320000, 8, 57, "SUCCESS", nil, "BID-57-2"},
		{340000, 9, 57, "SUCCESS", nil, "BID-57-3"},
		{380000, 7, 57, "SUCCESS", nil, "BID-57-4"},

		// Product 58: Serum Vitamin C (2 bids)
		{200000, 7, 58, "SUCCESS", nil, "BID-58-1"},
		{250000, 8, 58, "SUCCESS", nil, "BID-58-2"},

		// Product 59: Son Dior (2 bids)
		{600000, 8, 59, "SUCCESS", nil, "BID-59-1"},
		{700000, 9, 59, "SUCCESS", nil, "BID-59-2"},

		// Product 60: Innisfree Set (2 bids)
		{800000, 9, 60, "SUCCESS", nil, "BID-60-1"},
		{900000, 7, 60, "SUCCESS", nil, "BID-60-2"},
	}

	for i, bid := range biddingHistory {
		createdAt := now.Add(-time.Duration(len(biddingHistory)-i) * 15 * time.Minute)
		_, err := db.ExecContext(ctx, `
			INSERT INTO bidding_history (amount, bidder_id, product_id, status, reason, request_id, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, bid.amount, bid.bidderID, bid.productID, bid.status, bid.reason, bid.requestID, createdAt)
		if err != nil {
			return fmt.Errorf("error inserting bidding history: %v", err)
		}
	}

	log.Println("Bidding history seeded successfully")
	return nil
}

func seedComments(ctx context.Context, db *sql.DB) error {
	log.Println("Seeding comments...")

	now := time.Now()

	comments := []struct {
		productID int64
		senderID  int64
		content   string
	}{
		// Product 1: iPhone 15 Pro Max
		{1, 7, "Máy còn bảo hành bao lâu ạ?"},
		{1, 10, "Máy còn bảo hành 11 tháng chính hãng Apple Việt Nam ạ."},
		{1, 8, "Shop có giao hàng tận nơi không?"},
		{1, 10, "Có ạ, shop hỗ trợ giao hàng toàn quốc."},
		{1, 9, "Pin còn mấy phần trăm vậy shop?"},
		{1, 10, "Pin health 98% ạ, còn rất mới."},

		// Product 2: Samsung S24 Ultra
		{2, 8, "Máy có kèm ốp lưng không ạ?"},
		{2, 10, "Có kèm ốp lưng Samsung chính hãng ạ."},
		{2, 9, "Màn hình có dán kính cường lực chưa?"},
		{2, 10, "Dán rồi ạ, kính UV chống lóa."},

		// Product 5: MacBook Pro M3
		{5, 8, "MacBook này có kèm túi xách không shop?"},
		{5, 10, "Không kèm túi xách ạ, chỉ có phụ kiện theo máy thôi."},
		{5, 7, "Màn hình có vết xước không ạ?"},
		{5, 10, "Màn hình còn hoàn hảo, không vết xước."},

		// Product 6: Dell XPS 15
		{6, 9, "Laptop có nâng cấp RAM hay SSD không ạ?"},
		{6, 10, "Không ạ, máy còn zin 100%."},

		// Product 10: iPad Pro M2
		{10, 7, "iPad có kèm bút Apple Pencil không ạ?"},
		{10, 10, "Không kèm bút ạ, chỉ bán máy và phụ kiện đi kèm."},
		{10, 8, "Máy có bao da không shop?"},
		{10, 10, "Có bao da Smart Folio chính hãng Apple."},

		// Product 13: AirPods Pro 2
		{13, 9, "Pin sạc không dây được không ạ?"},
		{13, 11, "Được ạ, hộp sạc hỗ trợ MagSafe và Qi."},
		{13, 7, "Có kèm tips tai nghe không ạ?"},
		{13, 11, "Có đầy đủ 4 size tips ạ."},

		// Product 17: Apple Watch Series 9
		{17, 8, "Đồng hồ có chống nước không shop?"},
		{17, 11, "Có ạ, Apple Watch Series 9 chống nước chuẩn 50m."},
		{17, 9, "Còn bảo hành bao lâu?"},
		{17, 11, "Còn 9 tháng bảo hành Apple VN."},

		// Product 20: Áo Bomber Nam
		{20, 8, "Size L phù hợp bao nhiêu kg ạ?"},
		{20, 11, "Size L cho 65-75kg ạ."},

		// Product 28: Nike Air Force 1
		{28, 7, "Giày còn size 42 không shop?"},
		{28, 17, "Còn ạ, shop còn đủ size từ 40-44."},
		{28, 9, "Có kèm box không ạ?"},
		{28, 17, "Có đầy đủ box, giấy gói ạ."},

		// Product 33: Túi Michael Kors
		{33, 8, "Túi có kèm dây đeo không ạ?"},
		{33, 18, "Có 2 dây đeo: dây ngắn và dây dài ạ."},

		// Product 37: Bàn làm việc
		{37, 9, "Bàn có điều chỉnh độ cao không ạ?"},
		{37, 17, "Có ạ, điều chỉnh bằng điện, nhớ 4 vị trí."},

		// Product 40: Tủ lạnh Samsung
		{40, 8, "Tủ lạnh này tiêu thụ điện thế nào ạ?"},
		{40, 10, "Tủ này công nghệ Inverter nên tiết kiệm điện lắm ạ."},
		{40, 7, "Shop có lắp đặt miễn phí không?"},
		{40, 10, "Có ạ, miễn phí lắp đặt trong nội thành."},

		// Product 41: Máy giặt LG
		{41, 9, "Máy giặt có tiếng ồn nhiều không ạ?"},
		{41, 11, "Máy Inverter nên rất êm ạ."},

		// Product 46: Yamaha Exciter 155
		{46, 7, "Xe có đăng kiểm chưa ạ?"},
		{46, 15, "Đăng kiểm còn hạn 8 tháng ạ."},
		{46, 9, "Xe có tai nạn chưa ạ?"},
		{46, 15, "Chưa té ngã lần nào, còn zin nguyên."},

		// Product 47: SH Mode 2024
		{47, 8, "Xe còn bảo hành không ạ?"},
		{47, 15, "Còn 1 năm bảo hành Honda ạ."},

		// Product 49: Bộ sách Harry Potter
		{49, 7, "Sách có rách hay ố vàng không ạ?"},
		{49, 17, "Sách như mới, bảo quản cẩn thận."},

		// Product 51: Bóng đá Adidas
		{51, 9, "Bóng này có phù hợp sân cỏ nhân tạo không?"},
		{51, 15, "Có ạ, bóng chơi được cả sân cỏ tự nhiên và nhân tạo."},

		// Product 52: Vợt cầu lông Yonex
		{52, 8, "Vợt có căng cước chưa ạ?"},
		{52, 15, "Vợt đã căng BG65, lực 25lbs ạ."},

		// Product 53: Xe đạp Giant
		{53, 7, "Xe có kèm phụ kiện gì không ạ?"},
		{53, 15, "Có đèn, khóa, bình nước và giá để bình."},

		// Product 14: Sony WH-1000XM5
		{14, 7, "Tai nghe có chống ồn chủ động không?"},
		{14, 11, "Có ANC rất tốt, giảm 95% tiếng ồn."},

		// Product 29: Adidas Ultraboost 22
		{29, 8, "Giày có đế boost thật không ạ?"},
		{29, 17, "Đế boost chính hãng Adidas ạ."},

		// Product 35: Đồng hồ G-Shock
		{35, 9, "Đồng hồ có chống nước không ạ?"},
		{35, 17, "Chống nước 200m, bơi lặn được."},

		// Product 38: Ghế gaming E-Dra
		{38, 7, "Ghế có tựa đầu và tựa lưng không ạ?"},
		{38, 17, "Có cả 2, điều chỉnh được ạ."},

		// Product 43: Nồi cơm Zojirushi
		{43, 8, "Nồi có chế độ hẹn giờ không ạ?"},
		{43, 10, "Có đầy đủ chức năng hẹn giờ, giữ ấm."},

		// Product 12: Samsung Galaxy Tab S9+
		{12, 9, "Tab có kèm bút S Pen không ạ?"},
		{12, 10, "Có kèm S Pen chính hãng."},

		// Product 18: Galaxy Watch 6 Classic
		{18, 7, "Đồng hồ có đo nhịp tim không ạ?"},
		{18, 11, "Có ECG và đo SpO2 chính xác."},

		// Product 30: Jordan 1 Retro High
		{30, 8, "Giày có vàng đế không ạ?"},
		{30, 17, "Đế vẫn trắng, chưa bị vàng."},

		// Product 34: Balo Kanken Classic
		{34, 9, "Balo có nhiều ngăn không ạ?"},
		{34, 18, "Có 2 ngăn lớn và 2 ngăn nhỏ bên hông."},
	}

	for i, comment := range comments {
		createdAt := now.Add(-time.Duration(len(comments)-i) * 90 * time.Minute)
		_, err := db.ExecContext(ctx, `
			INSERT INTO comments (product_id, sender_id, content, created_at)
			VALUES ($1, $2, $3, $4)
		`, comment.productID, comment.senderID, comment.content, createdAt)
		if err != nil {
			return fmt.Errorf("error inserting comment: %v", err)
		}
	}

	log.Println("Comments seeded successfully")
	return nil
}

func seedOrders(ctx context.Context, db *sql.DB) error {
	log.Println("Seeding orders...")

	now := time.Now()

	orders := []struct {
		auctionID       int64
		winnerID        int64
		sellerID        int64
		finalPrice      float64
		status          string
		paymentMethod   *string
		paymentProof    *string
		shippingAddress *string
		shippingPhone   *string
		trackingNumber  *string
		shippingInvoice *string
		paidAt          *time.Time
		deliveredAt     *time.Time
		completedAt     *time.Time
		cancelledAt     *time.Time
		cancelReason    *string
	}{
		// Order 1: Honda Wave (Product 45)
		{
			auctionID:       45,
			winnerID:        7,
			sellerID:        17,
			finalPrice:      20000000,
			status:          "COMPLETED",
			paymentMethod:   strPtr("BANK_TRANSFER"),
			paymentProof:    strPtr("https://example.com/proof/payment-45.jpg"),
			shippingAddress: strPtr("123 Nguyễn Văn Linh, Quận 7, TP.HCM"),
			shippingPhone:   strPtr("0901234567"),
			trackingNumber:  strPtr("VN123456789"),
			shippingInvoice: strPtr("https://example.com/invoice/shipping-45.jpg"),
			paidAt:          timePtr(now.Add(-20 * 24 * time.Hour)),
			deliveredAt:     timePtr(now.Add(-15 * 24 * time.Hour)),
			completedAt:     timePtr(now.Add(-14 * 24 * time.Hour)),
			cancelledAt:     nil,
			cancelReason:    nil,
		},
		// Order 2: Sách Đắc Nhân Tâm (Product 48)
		{
			auctionID:       48,
			winnerID:        8,
			sellerID:        18,
			finalPrice:      80000,
			status:          "COMPLETED",
			paymentMethod:   strPtr("CASH"),
			paymentProof:    nil,
			shippingAddress: strPtr("456 Lê Văn Việt, Quận 9, TP.HCM"),
			shippingPhone:   strPtr("0912345678"),
			trackingNumber:  strPtr("VN987654321"),
			shippingInvoice: nil,
			paidAt:          timePtr(now.Add(-10 * 24 * time.Hour)),
			deliveredAt:     timePtr(now.Add(-7 * 24 * time.Hour)),
			completedAt:     timePtr(now.Add(-6 * 24 * time.Hour)),
			cancelledAt:     nil,
			cancelReason:    nil,
		},
		// Order 3: iPhone 14 Pro (Product 3) - Đang chờ thanh toán
		{
			auctionID:       3,
			winnerID:        9,
			sellerID:        10,
			finalPrice:      20000000,
			status:          "PENDING_PAYMENT",
			paymentMethod:   nil,
			paymentProof:    nil,
			shippingAddress: strPtr("789 Võ Văn Tần, Quận 3, TP.HCM"),
			shippingPhone:   strPtr("0923456789"),
			trackingNumber:  nil,
			shippingInvoice: nil,
			paidAt:          nil,
			deliveredAt:     nil,
			completedAt:     nil,
			cancelledAt:     nil,
			cancelReason:    nil,
		},
		// Order 4: Xiaomi 14 (Product 4) - Đang vận chuyển
		{
			auctionID:       4,
			winnerID:        7,
			sellerID:        11,
			finalPrice:      16000000,
			status:          "SHIPPING",
			paymentMethod:   strPtr("BANK_TRANSFER"),
			paymentProof:    strPtr("https://example.com/proof/payment-4.jpg"),
			shippingAddress: strPtr("321 Trần Hưng Đạo, Quận 1, TP.HCM"),
			shippingPhone:   strPtr("0934567890"),
			trackingNumber:  strPtr("VN111222333"),
			shippingInvoice: strPtr("https://example.com/invoice/shipping-4.jpg"),
			paidAt:          timePtr(now.Add(-5 * 24 * time.Hour)),
			deliveredAt:     nil,
			completedAt:     nil,
			cancelledAt:     nil,
			cancelReason:    nil,
		},
	}

	for _, order := range orders {
		_, err := db.ExecContext(ctx, `
			INSERT INTO orders (
				auction_id, winner_id, seller_id, final_price, status,
				payment_method, payment_proof, shipping_address, shipping_phone,
				tracking_number, shipping_invoice, paid_at, delivered_at,
				completed_at, cancelled_at, cancel_reason, created_at, updated_at
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		`, order.auctionID, order.winnerID, order.sellerID, order.finalPrice, order.status,
			order.paymentMethod, order.paymentProof, order.shippingAddress, order.shippingPhone,
			order.trackingNumber, order.shippingInvoice, order.paidAt, order.deliveredAt,
			order.completedAt, order.cancelledAt, order.cancelReason, now, now)
		if err != nil {
			return fmt.Errorf("error inserting order: %v", err)
		}
	}

	log.Println("Orders seeded successfully")
	return nil
}

func seedOrderMessages(ctx context.Context, db *sql.DB) error {
	log.Println("Seeding order messages...")

	now := time.Now()

	messages := []struct {
		orderID  int64
		senderID int64
		message  string
	}{
		// Order 1: Honda Wave - Giao dịch hoàn thành
		{1, 7, "Chào shop, em đã chuyển khoản 20 triệu rồi ạ."},
		{1, 17, "Shop đã nhận được thanh toán, sẽ gửi hàng trong ngày hôm nay."},
		{1, 7, "Cảm ơn shop ạ!"},
		{1, 17, "Hàng đã gửi, mã vận đơn: VN123456789"},
		{1, 7, "Em đã nhận được xe, xe rất đẹp. Cảm ơn shop nhiều!"},
		{1, 17, "Cảm ơn bạn đã tin tưởng, chúc bạn sử dụng vui vẻ!"},

		// Order 2: Sách Đắc Nhân Tâm - Giao dịch trực tiếp
		{2, 8, "Shop ơi, em muốn gặp trực tiếp để nhận hàng được không?"},
		{2, 18, "Được ạ, shop ở quận 1, bạn có thể đến lấy."},
		{2, 8, "Vâng, em sẽ qua chiều nay."},
		{2, 18, "Ok, hẹn bạn nhé!"},
		{2, 8, "Em đã nhận sách rồi ạ, cảm ơn shop!"},

		// Order 3: iPhone 14 Pro - Chờ thanh toán
		{3, 9, "Shop ơi, em chuyển khoản ngân hàng nào ạ?"},
		{3, 10, "Em chuyển vào STK Techcombank: 19036768686868, tên Nguyễn Văn A ạ."},
		{3, 9, "Em sẽ chuyển trong hôm nay ạ."},
		{3, 10, "Dạ, shop đợi bạn nha."},

		// Order 4: Xiaomi 14 - Đang vận chuyển
		{4, 7, "Em đã chuyển khoản rồi shop ơi."},
		{4, 11, "Shop đã nhận, đóng gói gửi hàng ngay."},
		{4, 11, "Hàng đã gửi, mã VN111222333. Dự kiến 2-3 ngày sẽ nhận được hàng."},
		{4, 7, "Cảm ơn shop ạ!"},
	}

	for i, msg := range messages {
		createdAt := now.Add(-time.Duration(len(messages)-i) * 6 * time.Hour)
		_, err := db.ExecContext(ctx, `
			INSERT INTO order_messages (order_id, sender_id, message, created_at)
			VALUES ($1, $2, $3, $4)
		`, msg.orderID, msg.senderID, msg.message, createdAt)
		if err != nil {
			return fmt.Errorf("error inserting order message: %v", err)
		}
	}

	log.Println("Order messages seeded successfully")
	return nil
}

func seedOrderRatings(ctx context.Context, db *sql.DB) error {
	log.Println("Seeding order ratings...")

	now := time.Now()

	ratings := []struct {
		orderID       int64
		buyerRating   *int
		buyerComment  *string
		sellerRating  *int
		sellerComment *string
	}{
		{
			orderID:       1,
			buyerRating:   intPtr(1),
			buyerComment:  strPtr("Người mua rất uy tín, thanh toán nhanh!"),
			sellerRating:  intPtr(1),
			sellerComment: strPtr("Shop giao hàng nhanh, xe đẹp như mô tả. Rất hài lòng!"),
		},
		{
			orderID:       2,
			buyerRating:   intPtr(1),
			buyerComment:  strPtr("Người mua dễ thương, lịch sự!"),
			sellerRating:  intPtr(1),
			sellerComment: strPtr("Sách đẹp, shop nhiệt tình!"),
		},
	}

	for _, rating := range ratings {
		_, err := db.ExecContext(ctx, `
			INSERT INTO order_ratings (
				order_id, buyer_rating, buyer_comment, seller_rating, seller_comment,
				buyer_rated_at, seller_rated_at, created_at, updated_at
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`, rating.orderID, rating.buyerRating, rating.buyerComment, rating.sellerRating, rating.sellerComment,
			now, now, now, now)
		if err != nil {
			return fmt.Errorf("error inserting order rating: %v", err)
		}
	}

	log.Println("Order ratings seeded successfully")
	return nil
}

func seedUserUpgradeRequests(ctx context.Context, db *sql.DB) error {
	log.Println("Seeding user upgrade requests...")

	now := time.Now()

	requests := []struct {
		userID          int64
		reason          string
		status          string
		rejectionReason *string
		createdAt       time.Time
		reviewedAt      *time.Time
		reviewedByID    *int64
	}{
		{
			userID:          9,
			reason:          "Tôi muốn bán các sản phẩm điện tử cũ của mình. Tôi có nhiều kinh nghiệm giao dịch online.",
			status:          "PENDING",
			rejectionReason: nil,
			createdAt:       now.Add(-24 * time.Hour),
			reviewedAt:      nil,
			reviewedByID:    nil,
		},
		{
			userID:          21,
			reason:          "Tôi có nguồn hàng ổn định và muốn kinh doanh trên sàn đấu giá.",
			status:          "PENDING",
			rejectionReason: nil,
			createdAt:       now.Add(-48 * time.Hour),
			reviewedAt:      nil,
			reviewedByID:    nil,
		},
		{
			userID:          23,
			reason:          "Tôi muốn thử bán một số đồ cũ không dùng nữa.",
			status:          "REJECTED",
			rejectionReason: strPtr("Tài khoản chưa đủ uy tín, vui lòng tham gia đấu giá thêm."),
			createdAt:       now.Add(-72 * time.Hour),
			reviewedAt:      timePtr(now.Add(-60 * time.Hour)),
			reviewedByID:    int64Ptr(15),
		},
	}

	for _, req := range requests {
		_, err := db.ExecContext(ctx, `
			INSERT INTO user_upgrade_requests (
				user_id, reason, status, rejection_reason, created_at, reviewed_at, reviewed_by_admin_id
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, req.userID, req.reason, req.status, req.rejectionReason, req.createdAt, req.reviewedAt, req.reviewedByID)
		if err != nil {
			return fmt.Errorf("error inserting user upgrade request: %v", err)
		}
	}

	log.Println("User upgrade requests seeded successfully")
	return nil
}

// Helper functions
func floatPtr(f float64) *float64 {
	return &f
}

func int64Ptr(i int64) *int64 {
	return &i
}

func intPtr(i int) *int {
	return &i
}

func strPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}
