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
			InsecureSkipVerify: true,
			ServerName:         cfg.DBHost,
		},
	})

	ctx := context.Background()
	if err := db.Ping(ctx); err != nil {
		log.Fatalf("Không thể kết nối database: %v", err)
	}

	log.Println("Kết nối PostgreSQL thành công!")
	return db
}
