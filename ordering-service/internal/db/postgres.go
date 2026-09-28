package db

import (
	"fmt"
	"log"
	"ordering-service/configs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Vientiane",
		configs.Cfg.Database.Host,
		configs.Cfg.Database.User,
		configs.Cfg.Database.Password,
		configs.Cfg.Database.DBName,
		configs.Cfg.Database.Port,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	log.Println("✅ Connected to PostgreSQL successfully (QR-ORDER)")
}
