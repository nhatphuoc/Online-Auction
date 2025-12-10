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

	// Create auto_bids table
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS auto_bids (
			id BIGSERIAL PRIMARY KEY,
			product_id BIGINT NOT NULL,
			bidder_id BIGINT NOT NULL,
			max_amount DOUBLE PRECISION NOT NULL,
			current_amount DOUBLE PRECISION NOT NULL DEFAULT 0,
			status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT auto_bids_status_check CHECK (status IN ('ACTIVE', 'WON', 'OUTBID', 'CANCELLED', 'EXPIRED'))
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating auto_bids table: %v", err)
	}

	// Create index on product_id for faster lookup
	_, err = db.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS idx_auto_bids_product_id ON auto_bids(product_id)
	`)
	if err != nil {
		return fmt.Errorf("error creating index on auto_bids: %v", err)
	}

	// Create index on bidder_id
	_, err = db.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS idx_auto_bids_bidder_id ON auto_bids(bidder_id)
	`)
	if err != nil {
		return fmt.Errorf("error creating index on auto_bids: %v", err)
	}

	// Create composite index on product_id and status for filtering active auto-bids
	_, err = db.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS idx_auto_bids_product_status ON auto_bids(product_id, status)
	`)
	if err != nil {
		return fmt.Errorf("error creating composite index on auto_bids: %v", err)
	}

	// Create index on max_amount for sorting
	_, err = db.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS idx_auto_bids_max_amount ON auto_bids(max_amount DESC)
	`)
	if err != nil {
		return fmt.Errorf("error creating index on max_amount: %v", err)
	}

	log.Println("Database schema initialized successfully!")
	return nil
}
