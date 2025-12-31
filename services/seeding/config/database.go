package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func ConnectDB(cfg *Config) *sql.DB {
	// Build connection string for PostgreSQL with SSL mode
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Không thể mở kết nối database: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Không thể kết nối database: %v", err)
	}

	log.Println("Kết nối database thành công!")
	return db
}
