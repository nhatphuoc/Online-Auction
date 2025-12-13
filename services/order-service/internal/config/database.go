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

	// Create orders table
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS orders (
			id BIGSERIAL PRIMARY KEY,
			auction_id BIGINT NOT NULL,
			winner_id BIGINT NOT NULL,
			seller_id BIGINT NOT NULL,
			final_price DOUBLE PRECISION NOT NULL,
			status VARCHAR(50) NOT NULL,
			payment_method VARCHAR(100),
			payment_proof TEXT,
			shipping_address TEXT,
			shipping_phone VARCHAR(20),
			tracking_number VARCHAR(100),
			shipping_invoice TEXT,
			delivered_at TIMESTAMP,
			completed_at TIMESTAMP,
			cancelled_at TIMESTAMP,
			cancel_reason TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT orders_status_check CHECK (status IN ('PENDING_PAYMENT', 'PAYMENT_CONFIRMED', 'ADDRESS_PROVIDED', 'INVOICE_SENT', 'DELIVERED', 'COMPLETED', 'CANCELLED'))
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating orders table: %v", err)
	}

	// Create indexes on orders
	_, err = db.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS idx_orders_auction_id ON orders(auction_id)
	`)
	if err != nil {
		return fmt.Errorf("error creating index on orders: %v", err)
	}

	_, err = db.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS idx_orders_winner_id ON orders(winner_id)
	`)
	if err != nil {
		return fmt.Errorf("error creating index on orders: %v", err)
	}

	_, err = db.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS idx_orders_seller_id ON orders(seller_id)
	`)
	if err != nil {
		return fmt.Errorf("error creating index on orders: %v", err)
	}

	_, err = db.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status)
	`)
	if err != nil {
		return fmt.Errorf("error creating index on orders: %v", err)
	}

	// Create order_messages table
	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS order_messages (
			id BIGSERIAL PRIMARY KEY,
			order_id BIGINT NOT NULL,
			sender_id BIGINT NOT NULL,
			message TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating order_messages table: %v", err)
	}

	// Create index on order_messages
	_, err = db.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS idx_order_messages_order_id ON order_messages(order_id)
	`)
	if err != nil {
		return fmt.Errorf("error creating index on order_messages: %v", err)
	}

	// Create order_ratings table
	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS order_ratings (
			id BIGSERIAL PRIMARY KEY,
			order_id BIGINT NOT NULL UNIQUE,
			buyer_rating INT CHECK (buyer_rating IN (-1, 1)),
			buyer_comment TEXT,
			seller_rating INT CHECK (seller_rating IN (-1, 1)),
			seller_comment TEXT,
			buyer_rated_at TIMESTAMP,
			seller_rated_at TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating order_ratings table: %v", err)
	}

	// Create index on order_ratings
	_, err = db.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS idx_order_ratings_order_id ON order_ratings(order_id)
	`)
	if err != nil {
		return fmt.Errorf("error creating index on order_ratings: %v", err)
	}

	log.Println("Database schema initialized successfully!")
	return nil
}
