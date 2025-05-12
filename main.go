package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Todo モデル定義
type Todo struct {
	gorm.Model
	Content string `json:"content"`
}

// DBConfig 環境変数からデータベース設定を取得
type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	Database string
}

// getDBConfig 環境変数からDB設定を取得
func getDBConfig() DBConfig {
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatalf("Invalid DB_PORT value: %v", err)
	}

	config := DBConfig{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		Database: os.Getenv("MYSQL_DATABASE"), // ここをMYSQL_DATABASEに変更
	}

	// デバッグ用：設定確認
	fmt.Printf("DB Config: %+v\n", config)
	return config
}

// connectionDB データベース接続
func connectionDB() (*gorm.DB, error) {
	config := getDBConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		config.User, config.Password, config.Host, config.Port, config.Database)

	// デバッグ用：接続情報確認
	fmt.Printf("Connecting to database at %s:%d/%s\n", config.Host, config.Port, config.Database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return db, err
}

func main() {
	r := gin.Default()
	db, err := connectionDB()

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// スキーマの自動マイグレーション
	err = db.AutoMigrate(&Todo{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Pingエンドポイント
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	fmt.Println("Database connection and setup successful")
	r.Run(":8080")
}
