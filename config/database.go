package config

import (
	"gorm.io/driver/postgres"
    "gorm.io/gorm"
    "log"
)

var DB *gorm.DB

func ConnectDB() {
    var err error
    DB, err = gorm.Open(postgres.Open(DBUrl), &gorm.Config{})
    if err != nil {
        log.Fatalf("❌ Failed to connect to database: %v", err)
    }

    log.Println("✅ Connected to PostgreSQL")
}