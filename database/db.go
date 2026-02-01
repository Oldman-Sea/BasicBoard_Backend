package database

import (
	"board/models"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
    var err error

    // .env 파일 로드
    godotenv.Load()

    // MySQL 연결 정보
    dbUser := getEnv("DB_USER", "root")
    dbPassword := getEnv("DB_PASSWORD", "password")
    dbHost := getEnv("DB_HOST", "127.0.0.1")
    dbPort := getEnv("DB_PORT", "3306")
    dbName := getEnv("DB_NAME", "basicboard_db")

    // DSN 생성
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        dbUser, dbPassword, dbHost, dbPort, dbName)

    // MySQL 연결
    DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // SQL 쿼리 로깅 활성화 (개발용)
    DB = DB.Debug()

    // 테이블 자동 생성 (개발용)
    err = DB.AutoMigrate(&models.Post{}, &models.SearchHistory{})
    if err != nil {
        log.Fatal("Failed to migrate:", err)
    }

    log.Println("Database connected!")
}

func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}