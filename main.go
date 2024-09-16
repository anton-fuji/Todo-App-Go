package main

import (
	"github.com/gin-gonic/gin"
)

// Todo構造体を定義
type TodoItem struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

// TodoItemsを格納する配列
var todoItems = []TodoItem{}

func main() {
	route := gin.Default()

	// TODOの取得
	route.GET("/todos", func(c *gin.Context) {
		c.IndentedJSON(200, todoItems)
	})

	// TODOの作成
	route.POST("/todos", func(c *gin.Context) {
		var newItem TodoItem

		// リクエストのJSONをTodoItemにバインド
		if err := c.BindJSON(&newItem); err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON data"})
			return
		}

		// 新しいTodoItemを追加
		todoItems = append(todoItems, newItem)
		c.IndentedJSON(201, newItem)
	})

	// ポート8080でサーバを実行
	route.Run(":8080")
}
