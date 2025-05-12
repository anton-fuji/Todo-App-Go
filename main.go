package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
		Database: os.Getenv("MYSQL_DATABASE"),
	}

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

func errorDB(db *gorm.DB, c *gin.Context) bool {
	if db.Error != nil {
		log.Printf("Error todos: %v", db.Error)
		c.AbortWithStatus(http.StatusInternalServerError)
		return true
	}
	return false
}

func listeners(r *gin.Engine, db *gorm.DB) {
	r.GET("/todo/delete", func(c *gin.Context) {
		id, _ := c.GetQuery("id")
		result := db.Delete(&Todo{}, id)
		if errorDB(result, c) {
			return
		}
	})

	r.POST("/todo/update", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.PostForm("id"))
		content := c.PostForm("content")
		var todo Todo
		result := db.Where("id = ?", id).Take(&todo)
		if errorDB(result, c) {
			return
		}
		todo.Content = content
		result = db.Save(&todo)
		if errorDB(result, c) {
			return
		}
	})

	r.POST("/todo/create", func(c *gin.Context) {
		content := c.PostForm("content")
		fmt.Println(c.Request.PostForm, content)
		result := db.Create(&Todo{Content: content})
		if errorDB(result, c) {
			return
		}
	})

	r.GET("/todo/list", func(c *gin.Context) {
		var todos []Todo
		result := db.Find(&todos)
		if errorDB(result, c) {
			return
		}
		fmt.Println(json.NewEncoder(os.Stdout).Encode(todos))
		c.JSON(http.StatusOK, todos)

		r.GET("/todo/get", func(c *gin.Context) {
			var todo Todo
			id, _ := c.GetQuery("id")
			result := db.First(&todo, id)
			if errorDB(result, c) {
				return
			}
			fmt.Println(json.NewEncoder(os.Stdout).Encode(todo))
			c.JSON(http.StatusOK, todo)
		})
	})
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

	listeners(r, db)

	fmt.Println("Database connection and setup successful")
	r.Run(":8080")
}
