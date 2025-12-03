package config

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"

	"github.com/go-pg/pg/v10"
)

func ConnectDB(cfg *Config) *pg.DB {
	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.DBHost, cfg.DBPort),
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		Database: cfg.DBName,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true, // Neon requires TLS but with relaxed verification
			ServerName:         cfg.DBHost,
		},
	})

	ctx := context.Background()
	if err := db.Ping(ctx); err != nil {
		log.Fatalf("Không thể kết nối database: %v", err)
	}

	log.Println("Kết nối database thành công!")
	return db
}

func InitSchema(db *pg.DB) error {
	ctx := context.Background()

	// Tạo bảng products (basic product info for comment service)
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS products (
			id SERIAL PRIMARY KEY,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("lỗi tạo bảng products: %v", err)
	}

	// Tạo bảng comments
	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS comments (
			id SERIAL PRIMARY KEY,
			product_id INTEGER NOT NULL,
			sender_id BIGINT NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return fmt.Errorf("lỗi tạo bảng comments: %v", err)
	}

	// Tạo index cho comments
	_, err = db.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS idx_comments_product ON comments(product_id);
		CREATE INDEX IF NOT EXISTS idx_comments_sender ON comments(sender_id);
		CREATE INDEX IF NOT EXISTS idx_comments_created_at ON comments(created_at DESC);
	`)
	if err != nil {
		return fmt.Errorf("lỗi tạo index cho comments: %v", err)
	}

	log.Println("Schema initialized successfully!")
	return nil
}
