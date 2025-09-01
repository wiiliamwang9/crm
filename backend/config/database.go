package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "host=db.lamdar.cn user=postgres password=tpg1688 dbname=crm port=9524 sslmode=disable TimeZone=Asia/Shanghai client_encoding=UTF8"

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Database connected successfully")
	DB = database
}
