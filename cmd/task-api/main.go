package main

import (
	"log"
	"task-api/internal/item"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// POST 	/items
// GET 		/items?status=xxxxx
// PATCH	/items/:id

// GET 		/items/:id
// PUT		/items/:id
// DELETE 	/items/:id

func main() {
	// Connect database
	db, err := gorm.Open(
		postgres.Open(
			"postgres://postgres:password@localhost:5432/task",
		),
	)
	if err != nil {
		log.Panic(err)
	}

	// Controller
	controller := item.NewController(db)

	// Router
	r := gin.Default()

	config := cors.DefaultConfig()
	// frontend URL
	config.AllowOrigins = []string{
		"http://127.0.0.1:8000",
	}
	r.Use(cors.New(config))

	// Register router
	r.POST("/items", controller.CreateItem)
	r.GET("/items", controller.FindItems)
	r.PATCH("/items/:id", controller.UpdateItemStatus)

	// Start server
	if err := r.Run(); err != nil {
		log.Panic(err)
	}
}
