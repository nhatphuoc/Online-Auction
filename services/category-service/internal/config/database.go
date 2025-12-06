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

	// Create categories table
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS categories (
			id BIGSERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			slug VARCHAR(255) NOT NULL UNIQUE,
			description TEXT,
			parent_id BIGINT,
			level INT NOT NULL DEFAULT 1,
			is_active BOOLEAN NOT NULL DEFAULT true,
			display_order INT DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (parent_id) REFERENCES categories(id) ON DELETE SET NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating categories table: %v", err)
	}

	// Create index on parent_id
	_, err = db.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS idx_categories_parent_id ON categories(parent_id)
	`)
	if err != nil {
		return fmt.Errorf("error creating index on categories: %v", err)
	}

	// Create index on level
	_, err = db.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS idx_categories_level ON categories(level)
	`)
	if err != nil {
		return fmt.Errorf("error creating index on categories: %v", err)
	}

	// Create products table
	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS products (
			id BIGSERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			category_id BIGINT NOT NULL,
			seller_id BIGINT NOT NULL,
			starting_price DOUBLE PRECISION NOT NULL,
			current_price DOUBLE PRECISION,
			buy_now_price DOUBLE PRECISION,
			step_price DOUBLE PRECISION NOT NULL,
			status VARCHAR(255) NOT NULL,
			thumbnail_url TEXT,
			auto_extend BOOLEAN NOT NULL DEFAULT false,
			end_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE RESTRICT,
			CONSTRAINT products_status_check CHECK (status IN ('ACTIVE', 'FINISHED', 'PENDING', 'REJECTED'))
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating products table: %v", err)
	}

	// Create index on category_id
	_, err = db.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS idx_products_category_id ON products(category_id)
	`)
	if err != nil {
		return fmt.Errorf("error creating index on products: %v", err)
	}

	// Create index on status
	_, err = db.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS idx_products_status ON products(status)
	`)
	if err != nil {
		return fmt.Errorf("error creating index on products: %v", err)
	}

	// Create product_images table
	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS product_images (
			product_id BIGINT NOT NULL,
			image_url VARCHAR(255) NOT NULL,
			FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating product_images table: %v", err)
	}

	// Create index on product_id
	_, err = db.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS idx_product_images_product_id ON product_images(product_id)
	`)
	if err != nil {
		return fmt.Errorf("error creating index on product_images: %v", err)
	}

	log.Println("Database schema initialized successfully!")
	return nil
}
